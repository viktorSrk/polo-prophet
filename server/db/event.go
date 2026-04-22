package db

import "database/sql"

type Event struct {
	ID     int
	GameID int
	Team   string
	Player string
	Type   string
}

func CreateEvent(db *sql.DB, event Event) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO events
		(game_id, team, player, type)
		VALUES (?, ?, ?, ?)
	`, event.GameID, event.Team, event.Player, event.Type)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func DeleteEventsForGame(db *sql.DB, game_id int64) error {
	_, err := db.Exec(`
		DELETE FROM events
		WHERE game_id = ?
	`, game_id)
	return err
}

func EventsExists(db *sql.DB, game_id int) (bool, error) {
	var id int64
	err := db.QueryRow(`
		SELECT id FROM events
		WHERE game_id = ?
	`, game_id).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
