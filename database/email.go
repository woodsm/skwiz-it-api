package database

import (
	"log"
)

func GetEmailAddresses(drawingId int64) []string {
	var db = getDatabase()
	defer db.Close()

	rows, err := db.Query(`
							SELECT
							  DISTINCT (email) AS email
							FROM section AS s
							  INNER JOIN app_user AS a ON s.app_user_id = a.id
							WHERE s.drawing_id = ?`)
	defer rows.Close()

	if err != nil {
		log.Fatalf("Unable to get all emails for drawing %d : %q\n", drawingId, err)
	}

	var results []string

	for rows.Next() {

		email := ""

		err := rows.Scan(&email)

		if err != nil {
			log.Print(err)
		}

		results = append(results, email)
	}
	return results
}
