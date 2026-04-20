package scrape

import (
	"database/sql"
	"log"
	"polo-prophet/server/db"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

var database *sql.DB

func ScrapeLeagues(db *sql.DB) {
	database = db

	c := colly.NewCollector()

	url := "https://www.wasserball-team-deutschland.de/"

	c.OnHTML(".league-row", getLeagueInfo)

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}

func getLeagueInfo(e *colly.HTMLElement) {
	name := e.ChildText(".league-name")
	gender := e.ChildText(".muted.small")
	subdomain := e.ChildAttr(".league-name", "href")

	id, err := db.LeagueExists(database, name)
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		id, err = db.CreateLeague(database, db.League{ID: 0, Name: name, Gender: gender})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := db.UpdateLeague(database, int(id), db.League{ID: 0, Name: name, Gender: gender})
		if err != nil {
			log.Fatal(err)
		}
	}

	season_start, ltype, scrape_id, group := getLeagueScrapeInfo(subdomain)
	league_id := id

	id, err = db.LeagueScrapeInfoExists(database, league_id)
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		_, err = db.CreateLeagueScrapeInfo(database, db.LeagueScrapeInfo{
			ID:          0,
			LeagueID:    int(league_id),
			SeasonStart: season_start,
			Type:        ltype,
			ScrapeID:    scrape_id,
			Group:       group,
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := db.UpdateLeagueScrapeInfo(database, int(id), db.LeagueScrapeInfo{
			ID:          0,
			LeagueID:    int(league_id),
			SeasonStart: season_start,
			Type:        ltype,
			ScrapeID:    scrape_id,
			Group:       group,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	// fmt.Println(getLeagueScrapeInfo(subdomain))
}

func getLeagueScrapeInfo(subdomain string) (int, string, int, string) {
	splits := strings.Split(subdomain, "/")

	season_start, _ := strconv.Atoi(splits[2])
	ltype := splits[3]

	splits = strings.Split(splits[4], "?")
	scrape_id, _ := strconv.Atoi(splits[0])

	group := ""
	if len(splits) > 1 {
		splits = strings.Split(splits[1], "=")
		group = splits[1]
	}

	return season_start, ltype, scrape_id, group
}
