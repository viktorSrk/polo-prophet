package db

import "database/sql"

type League struct {
	ID		int
	Name	string
	Gender	string
}

type LeagueScrapeInfo struct {
	ID			int
	LeagueID	int
	SeasonStart	int
	Type		string
	ScrapeID	int
	Group		string
}

func CreateLeague(db *sql.DB, league League) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO leagues
		(name, gender)
		VALUES (?, ?)
	`, league.Name, league.Gender)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateLeague(db *sql.DB, id int, league League) error {
	_, err := db.Exec("UPDATE leagues SET name = ?, gender = ? WHERE id = ?", league.Name, league.Gender, id)
	return err
}

func LeagueExists(db *sql.DB, name string) (int64, error) {
	var id int64
	err := db.QueryRow(`
		SELECT id FROM leagues
		WHERE name = ?
	`, name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func CreateLeagueScrapeInfo(db *sql.DB, info LeagueScrapeInfo) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO leagues_scrapeInfo
		(league_id, season_start, type, scrape_id, 'group')
		VALUES (?, ?, ?, ?, ?)
	`, info.LeagueID, info.SeasonStart, info.Type, info.ScrapeID, info.Group)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateLeagueScrapeInfo(db *sql.DB, id int, info LeagueScrapeInfo) error {
	_, err := db.Exec("UPDATE leagues_scrapeInfo SET league_id = ?, season_start = ?, type = ?, scrape_id = ?, 'group' = ? WHERE id = ?", info.LeagueID, info.SeasonStart, info.Type, info.ScrapeID, info.Group, id)
	return err
}

func LeagueScrapeInfoExists(db *sql.DB, league_id int64) (int64, error) {
	var id int64
	err := db.QueryRow(`
		SELECT id FROM leagues_scrapeInfo
		WHERE league_id = ?
	`, league_id).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
