package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/benkauffman/skwiz-it-api/model"
	"github.com/benkauffman/skwiz-it-api/notification"
	"github.com/benkauffman/skwiz-it-api/helper"
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
		setDrawingAsComplete(drawingId)
		notification.SendEmail(drawingId)
	} else {
		log.Printf("Drawing %d only has %d/%d parts supplied...", drawingId, qty, len(helper.GetSections()))
	}
}

func setDrawingAsComplete(drawingId int64) {
	var db = getDatabase()
	defer db.Close()

	_, err := db.Exec(`UPDATE drawing SET completed = NOW() WHERE id = ?`, drawingId)

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

	rows, err := db.Query(selectSingle, id)
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

func GetDrawings() []model.Drawing {
	var db = getDatabase()
	defer db.Close()

	rows, err := db.Query(selectAll)
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

var selectSingle = `
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
WHERE
  d.id = ?
GROUP BY
  d.id,
  d.url`

var selectAll = `
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
GROUP BY
  d.id,
  d.url
`
