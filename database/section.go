package database

import (
	"../model"
	"../helper"

	"log"
	"math/rand"
)

func SaveSection(userId int64, typeOf string, url string) (model.Drawing, error) {

	groupId := addToDrawing(userId, typeOf, url)

	drawing, err := GetDrawing(groupId)
	if err != nil {
		return drawing, err
	}

	return drawing, nil
}

func GetNeededSection() (string) {

	log.Printf("Looking for needed section...")

	var db = getDatabase()
	defer db.Close()

	sql := `
SELECT
  MIN(qty) AS qty,
  type     AS type
FROM (
       SELECT
         COUNT(drawing_id) AS qty,
         'top'             AS type
       FROM section
       WHERE type = 'top'

       UNION

       SELECT
         COUNT(drawing_id) AS qty,
         'middle'          AS type
       FROM section
       WHERE type = 'middle'

       UNION

       SELECT
         COUNT(drawing_id) AS qty,
         'bottom'          AS type
       FROM section
       WHERE type = 'bottom'
     ) AS smry
GROUP BY type
ORDER BY qty ASC
LIMIT 1
`
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		log.Fatalf("Unable get missing section from database : %q\n", err)
	}

	section := ""

	for rows.Next() {
		var (
			qty = 0
		)

		err = rows.Scan(&qty, &section)

		if err != nil {
			log.Print(err)
		}
	}

	if err != nil {
		log.Print(err)
	}

	if section == "" {
		sections := helper.GetSections()
		section = sections[rand.Intn(len(sections))]
		log.Printf("Wasn't able to find needed section, picking a random one... %s", section)
	} else {
		log.Printf("Section \"%s\" provided by the database...", section)
	}

	return section
}

func addToDrawing(userId int64, typeOf string, url string) int64 {
	var db = getDatabase()
	defer db.Close()

	sql := `INSERT INTO section 
			(drawing_id, type, app_user_id, url, created, updated) 
			VALUES (?, ?, ?, ?, NOW(), NOW())`

	drawingId := getMissingDrawingId(typeOf)

	existing := drawingId >= 1

	if !existing {
		log.Printf("Unable to find an existing drawing for section %s", typeOf)
		log.Printf("Creating new drawing for section %s", typeOf)
		drawingId = CreateDrawing()
	} else {
		log.Printf("Using existing drawing %d for section %s", drawingId, typeOf)
	}

	_, err := db.Exec(sql, drawingId, typeOf, userId, url)

	if err != nil {
		log.Fatalf("Unable to create %s section for drawing %d : %q\n", typeOf, drawingId, err)
	} else {
		log.Printf("Created section %s for drawing %d", typeOf, drawingId)
		if existing {
			log.Printf("Running background task to validate completion of drawing...")
			go IsDrawingComplete(drawingId)
		}
	}

	return drawingId
}

func getMissingDrawingId(typeOf string) (id int64) {
	log.Printf("Looking for drawing with missing %s section...", typeOf)

	var db = getDatabase()
	defer db.Close()

	sql := "SELECT id FROM drawing WHERE id NOT IN (SELECT drawing_id FROM section WHERE `type` = ?) LIMIT 1"
	row, err := db.QueryRow(sql, typeOf)

	if err != nil {
		log.Fatalf("Unable to check for missing parts in drawings : %q\n", err)
	}

	err = row.Scan(&id)
	if err != nil {
		log.Print(err)
	}

	return id
}
