package database

import (
	"github.com/benkauffman/skwiz-it-api/model"
	"github.com/benkauffman/skwiz-it-api/validation"
	"fmt"
)

func UpsertUser(user *model.User) (*model.User, error) {

	if !validation.IsValidEmail(user.Email) {
		return nil, fmt.Errorf("email address %s is not valid", user.Email)
	}

	var db = getDatabase()
	defer db.Close()

	sql := `INSERT INTO app_user (email, created, name, updated) VALUES (?, NOW(), ?, NOW())
			ON DUPLICATE KEY UPDATE updated = VALUES(updated)`

	res, err := db.Exec(sql, user.Email, user.Name)

	if err != nil {
		println("Exec err:", err.Error())
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			println("Error:", err.Error())
		} else {
			user.Id = id
			println("Upserted user : ", id)
		}
	}

	return user, nil
}
