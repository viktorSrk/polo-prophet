package db

import "database/sql"

type Game struct {
	ID         int
	LeagueID   int
	GameNumber string
	Team1      string
	Team2      string
	DateTime   int
}

type GameScrapeInfo struct {
	ID             int
	GameID         int
	SeasonStart    int
	Type           string
	LeagueScrapeID int
	GameScrapeID   int
	Group          string
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

func GetAllGameIds(db *sql.DB) ([]int64, error) {
	var ids []int64
	rows, err := db.Query(`
		SELECT id FROM games
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func CreateGameScrapeInfo(db *sql.DB, info GameScrapeInfo) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO games_scrapeInfo
		(game_id, season_start, type, league_scrape_id, game_scrape_id, 'group')
		VALUES (?, ?, ?, ?, ?, ?)
	`, info.GameID, info.SeasonStart, info.Type, info.LeagueScrapeID, info.GameScrapeID, info.Group)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateGameScrapeInfo(db *sql.DB, id int, info GameScrapeInfo) error {
	_, err := db.Exec(`
		UPDATE games_scrapeInfo
		SET
			game_id = ?,
			season_start = ?,
			type = ?,
			league_scrape_id = ?,
			game_scrape_id = ?,
			'group' = ?
		WHERE id = ?
	`, info.GameID, info.SeasonStart, info.Type, info.LeagueScrapeID, info.GameScrapeID, info.Group, id)
	return err
}

func GameScrapeInfoExists(db *sql.DB, game_id int64) (int64, error) {
	var id int64
	err := db.QueryRow(`
		SELECT id FROM games_scrapeInfo
		WHERE game_id = ?
	`, game_id).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetGameScrapeInfo(db *sql.DB, id int64) (GameScrapeInfo, error) {
	var did int
	var game_id int
	var season_start int
	var ltype string
	var league_scrape_id int
	var game_scrape_id int
	var group string
	err := db.QueryRow(`
		SELECT * FROM games_scrapeInfo
		WHERE id = ?
	`, id).Scan(&did, &game_id, &season_start, &ltype, &league_scrape_id, &game_scrape_id, &group)
	if err == sql.ErrNoRows {
		return GameScrapeInfo{}, nil
	}
	if err != nil {
		return GameScrapeInfo{}, err
	}
	return GameScrapeInfo{
		ID:             did,
		GameID:         game_id,
		SeasonStart:    season_start,
		Type:           ltype,
		LeagueScrapeID: league_scrape_id,
		GameScrapeID:   game_scrape_id,
		Group:          group,
	}, nil
}
