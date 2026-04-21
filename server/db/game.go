package db

import "database/sql"

type Game struct {
	ID int
	LeagueID int
	GameNumber string
	Team1 string
	Team2 string
	DateTime int
}

func CreateGame(db *sql.DB, game Game) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO games
		(league_id, game_number, team1, team2, datetime)
		VALUES (?, ?, ?, ?, ?)
	`, game.LeagueID, game.GameNumber, game.Team1, game.Team2, game.DateTime)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateGame(db *sql.DB, id int, game Game) error {
	_, err := db.Exec(`
		UPDATE games
		SET
			league_id = ?,
			game_number = ?,
			team1 = ?,
			team2 = ?,
			datetime = ?
		WHERE id = ?
	`, game.LeagueID, game.GameNumber, game.Team1, game.Team2, game.DateTime, id)
	return err
}

func GameExists(db *sql.DB, league_id int, game_number string) (int64, error) {
	var id int64
	err := db.QueryRow(`
		SELECT id FROM games
		WHERE league_id = ?
		AND game_number = ?
	`, league_id, game_number).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
