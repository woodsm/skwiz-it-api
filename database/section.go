package database

import (
	"log"

	"github.com/benkauffman/skwiz-it-api/model"
)

func SaveSection(userId int, typeOf string, url string) (model.Drawing, error) {

	groupId := addToDrawing(userId, typeOf, url)

	drawing, err := GetDrawing(groupId)
	if err != nil {
		return drawing, err
	}

	return drawing, nil
}

func addToDrawing(userId int, typeOf string, url string) int64 {
	var db = getDatabase()
	defer db.Close()

	sql := `INSERT INTO section 
			(drawing_id, type, app_user_id, url, created, updated) 
			VALUES (?, ?, ?, ?, NOW(), NOW())`

	drawingId := getMissingDrawingId(typeOf)

	if drawingId <= 0 {
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
