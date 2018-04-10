package database

import (
	"../config"

	"log"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql" // Imports the mysql adapter.
)

var conf = config.LoadConfig()

func getDatabase() sqlbuilder.Database {

	var settings = mysql.ConnectionURL{
		Database: conf.MySQL.Database,
		Host:     conf.MySQL.Host,
		User:     conf.MySQL.User,
		Password: conf.MySQL.Password,
	}

	// Attempting to establish a connection to the database.
	sess, err := mysql.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}

	return sess

}

func CheckHealth() (bool) {
	var db = getDatabase()
	defer db.Close()

	row, err := db.QueryRow("SELECT COUNT(id) AS qty FROM app_user")
	if err != nil {
		return false
	}

	qty := -1
	err = row.Scan(&qty)
	if err != nil {
		log.Print(err)
	}

	return qty >= 0

}
