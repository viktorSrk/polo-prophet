package scrape

import (
	"database/sql"
	"fmt"
	"log"
	"polo-prophet/server/db"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeGames(database *sql.DB, league_id int64) {
	c := colly.NewCollector()

	scrapeInfo_id, err := db.LeagueScrapeInfoExists(database, league_id)
	if err != nil {
		log.Fatal(err)
	}
	scrapeInfo, err := db.GetScrapeInfo(database, scrapeInfo_id)
	if err != nil {
		log.Fatal(err)
	}

	domain := "https://www.wasserball-team-deutschland.de" + createLeagueSubdomain(scrapeInfo)

	c.OnHTML(".match-card", func(e *colly.HTMLElement) {
		getGameInfo(e, database, int(league_id))
	})

	err = c.Visit(domain)
	if err != nil {
		log.Print(err)
	}
}

func getGameInfo(e *colly.HTMLElement, database *sql.DB, league_id int) {
	game_number := strings.TrimPrefix(e.ChildText("span.tag"), "Spiel ")

	teams := e.ChildTexts(".team-name")
	var team1 string
	var team2 string
	if len(teams) == 2 {
		team1 = teams[0]
		team2 = teams[1]
	}

	datetime := 0

	id, err := db.GameExists(database, league_id, game_number)
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		db.CreateGame(database, db.Game{
			ID:         0,
			LeagueID:   league_id,
			GameNumber: game_number,
			Team1:      team1,
			Team2:      team2,
			DateTime:   datetime,
		})
	} else {
		db.UpdateGame(database, int(id), db.Game{
			ID:         0,
			LeagueID:   league_id,
			GameNumber: game_number,
			Team1:      team1,
			Team2:      team2,
			DateTime:   datetime,
		})
	}

}

func createLeagueSubdomain(info db.LeagueScrapeInfo) string {
	ret := fmt.Sprintf("/league/%d/%s/%d", info.SeasonStart, info.Type, info.ScrapeID)
	if info.Group != "" {
		ret += fmt.Sprintf("?group=%s", info.Group)
	}
	return ret
}
