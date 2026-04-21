package scrape

import (
	"database/sql"
	"fmt"
	"log"
	"polo-prophet/server/db"

	"github.com/gocolly/colly/v2"
)

func ScrapeGames(database *sql.DB, league_id int64) {
	c := colly.NewCollector()

	id, err := db.LeagueScrapeInfoExists(database, league_id)
	if err != nil {
		log.Fatal(err)
	}
	scrapeInfo, err := db.GetScrapeInfo(database, id)
	if err != nil {
		log.Fatal(err)
	}

	domain := "https://www.wasserball-team-deutschland.de" + createLeagueSubdomain(scrapeInfo)
	fmt.Println(league_id)
	fmt.Println(scrapeInfo)
	fmt.Println(domain)

	c.OnHTML(".match-card", getGameInfo)

	err = c.Visit(domain)
	if err != nil {
		log.Print(err)
	}
}

func getGameInfo(e *colly.HTMLElement) {
	teams := e.ChildTexts(".team-name")
	if len(teams) == 2 {
		fmt.Println("Home: " + teams[0] + ", Away: " + teams[1])
	}
}

func createLeagueSubdomain(info db.LeagueScrapeInfo) string {
	ret := fmt.Sprintf("/league/%d/%s/%d", info.SeasonStart, info.Type, info.ScrapeID)
	if info.Group != "" {
		ret += fmt.Sprintf("?group=%s", info.Group)
	}
	return ret
}
