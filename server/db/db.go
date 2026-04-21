package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables(db)
	return db
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS leagues (
			id			INTEGER PRIMARY KEY AUTOINCREMENT,
			name		TEXT NOT NULL,
			gender		TEXT
		);
		CREATE TABLE IF NOT EXISTS leagues_scrapeInfo (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			league_id		INTEGER NOT NULL UNIQUE,
			season_start	INTEGER NOT NULL,
			type			STRING NOT NULL,
			scrape_id		INTEGER NOT NULL,
			'group'			STRING,
			FOREIGN KEY (league_id) REFERENCES leagues(id)
		);
		CREATE TABLE IF NOT EXISTS games (
			id 			INTEGER PRIMARY KEY AUTOINCREMENT,
			league_id 	INTEGER NOT NULL,
			game_number STRING NOT NULL,
			team1 		STRING,
			team2 		STRING,
			datetime 	INTEGER,
			FOREIGN KEY (league_id) REFERENCES leagues(id)
		);
		CREATE TABLE IF NOT EXISTS games_scrapeInfo (
			id					INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id				INTEGER NOT NULL UNIQUE,
			season_start		INTEGER NOT NULL,
			type				STRING NOT NULL,
			league_scrape_id 	INTEGER NOT NULL,
			game_scrape_id 		INTEGER NOT NULL,
			'group'				STRING,
			FOREIGN KEY (game_id) REFERENCES games(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
