package database

import (
	"../model"
	"../notification"
	"../helper"
	"../image"
	"../storage"

	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func IsDrawingComplete(drawingId int64) {
	log.Printf("Checking if drawing %d is complete and handling accordingly...", drawingId)
	var db = getDatabase()
	defer db.Close()

	row, err := db.QueryRow("SELECT COUNT(drawing_id) AS qty FROM section WHERE drawing_id = ?", drawingId)

	if err != nil {
		log.Fatalf("Unable to get section count : %q\n", err)
	}

	qty := 0
	err = row.Scan(&qty)
	if err != nil {
		log.Print(err)
	}

	if qty == len(helper.GetSections()) {
		log.Printf("Drawing %d has been completed... Sending out emails and setting as completed...", drawingId)

		//combine all the drawings into a full drawing
		drawing, err := GetDrawing(drawingId)
		helper.CheckError(err)
		b64 := image.CreateFullDrawing(drawing)

		//save the full drawing to S3
		fileId, err := storage.SaveToS3(b64)
		helper.CheckError(err)
		url := helper.GetUrl(fileId)

		setDrawingAsComplete(drawingId, url)
		sendEmails(drawingId)
	} else {
		log.Printf("Drawing %d only has %d/%d parts supplied...", drawingId, qty, len(helper.GetSections()))
	}
}

func sendEmails(drawingId int64) {
	log.Printf("Sending email to users for drawing %d completion", drawingId)
	for _, emailAddr := range GetEmailAddresses(drawingId) {
		log.Printf("Sending email in background thread to %s for drawing %d", emailAddr, drawingId)
		go notification.SendEmail(emailAddr, drawingId)
	}
}

func setDrawingAsComplete(drawingId int64, url string) {
	var db = getDatabase()
	defer db.Close()

	_, err := db.Exec(`UPDATE drawing SET completed = NOW(), url = ? WHERE id = ?`, url, drawingId)

	if err != nil {
		log.Fatalf("Unable to set drawing %d completed datetime : %q\n", drawingId, err)
	}
}

func CreateDrawing() int64 {
	var db = getDatabase()
	defer db.Close()

	res, err := db.Exec(`	INSERT INTO drawing (id, url, created, updated, completed) 
								VALUES (0, NULL, NOW(), NOW(), NULL)`)

	if err != nil {
		log.Fatalf("Unable to create drawing : %q\n", err)
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			log.Fatalf("Unable to get id for new drawing : %q\n", err)
		} else {
			log.Printf("Created new drawing %d", id)
			return id
		}
	}
	return 0
}

func GetDrawing(id int64) (model.Drawing, error) {
	var db = getDatabase()
	defer db.Close()

	rows, err := db.Query(buildSelectQuery(id, nil))
	defer rows.Close()

	if err != nil {
		log.Fatalf("Unable to get drawing : %q\n", err)
	}

	var list = parse(rows)

	if len(list) >= 1 {
		return list[0], nil
	}

	return model.Drawing{}, fmt.Errorf("not found: unable to find drawing with id %d", id)
}

func GetDrawings(user *model.User) []model.Drawing {
	var db = getDatabase()
	defer db.Close()

	rows, err := db.Query(buildSelectQuery(0, user))
	defer rows.Close()

	if err != nil {
		log.Fatalf("Unable to get all drawings : %q\n", err)
	}

	return parse(rows)
}

func parse(rows *sql.Rows) []model.Drawing {
	var results []model.Drawing

	for rows.Next() {

		d := model.Drawing{}

		err := rows.Scan(
			&d.Id,
			&d.Url,
			&d.Top.Url,
			&d.Top.Name,
			&d.Top.Email,
			&d.Middle.Url,
			&d.Middle.Name,
			&d.Middle.Email,
			&d.Bottom.Url,
			&d.Bottom.Name,
			&d.Bottom.Email,
		)

		if err != nil {
			log.Print(err)
		}

		results = append(results, d)
	}
	return results
}

func buildSelectQuery(drawingId int64, user *model.User) (string) {
	var s = `
SELECT
  d.id               AS drawing_id,
  d.url              AS drawing_url,

  MAX(CASE WHEN s.type = 'top'
    THEN s.url
      ELSE NULL END) AS section_top_url,

  MAX(CASE WHEN s.type = 'top'
    THEN a.name
      ELSE NULL END) AS section_top_name,

  MAX(CASE WHEN s.type = 'top'
    THEN a.email
      ELSE NULL END) AS section_top_email,

  MAX(CASE WHEN s.type = 'middle'
    THEN s.url
      ELSE NULL END) AS section_middle_url,

  MAX(CASE WHEN s.type = 'middle'
    THEN a.name
      ELSE NULL END) AS section_middle_name,

  MAX(CASE WHEN s.type = 'middle'
    THEN a.email
      ELSE NULL END) AS section_middle_email,

  MAX(CASE WHEN s.type = 'bottom'
    THEN s.url
      ELSE NULL END) AS section_bottom_url,

  MAX(CASE WHEN s.type = 'bottom'
    THEN a.name
      ELSE NULL END) AS section_bottom_name,

  MAX(CASE WHEN s.type = 'bottom'
    THEN a.email
      ELSE NULL END) AS section_bottom_email

FROM drawing AS d
  LEFT JOIN section s ON d.id = s.drawing_id
  LEFT JOIN app_user a ON s.app_user_id = a.id

WHERE d.id <> 0

`

	if drawingId >= 1 {
		s += " AND d.id = " + strconv.FormatInt(drawingId, 10)

	}

	if user != nil && user.Id >= 1 {
		s += " AND d.id IN ("
		s += " SELECT drawing_id FROM section WHERE app_user_id = " + strconv.FormatInt(user.Id, 10)
		s += ") "

	}

	s += `

GROUP BY
  d.id,
  d.url
`
	return s
}
