package scrape

import (
	"database/sql"
	"fmt"
	"log"
	"polo-prophet/server/db"
	"strconv"
	"strings"
	"time"

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

	c.OnError(func(r *colly.Response, err error) {
		// log.Printf("%d: Status: %d\n", game_id, r.StatusCode)
		switch r.StatusCode {
		case 429:
			time.Sleep(2000 * time.Millisecond) // fallback
			r.Request.Retry()
		case 502:
			// log.Println(domain)
		default:
			log.Println(err)
		}
	})

	c.Visit(domain)
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

	var datetime int
	datetime_string := e.ChildText("span.pill")
	t, err := time.Parse("02.01.06, 15:04 Uhr", datetime_string)
	if err != nil {
		datetime = 0
	} else {
		datetime = int(t.UnixMilli())
	}

	id, err := db.GameExists(database, league_id, game_number)
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		id, err = db.CreateGame(database, db.Game{
			ID:         0,
			LeagueID:   league_id,
			GameNumber: game_number,
			Team1:      team1,
			Team2:      team2,
			DateTime:   datetime,
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = db.UpdateGame(database, int(id), db.Game{
			ID:         0,
			LeagueID:   league_id,
			GameNumber: game_number,
			Team1:      team1,
			Team2:      team2,
			DateTime:   datetime,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	subdomain := e.ChildAttr("a.action", "href")
	season_start, ltype, league_scrape_id, game_scrape_id, group := getGameScrapeInfo(subdomain)
	game_id := id

	id, err = db.GameScrapeInfoExists(database, game_id)
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		_, err = db.CreateGameScrapeInfo(database, db.GameScrapeInfo{
			ID: 0,
			GameID: int(game_id),
			SeasonStart: season_start,
			Type: ltype,
			LeagueScrapeID: league_scrape_id,
			GameScrapeID: game_scrape_id,
			Group: group,
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = db.UpdateGameScrapeInfo(database, int(id), db.GameScrapeInfo{
			ID: 0,
			GameID: int(game_id),
			SeasonStart: season_start,
			Type: ltype,
			LeagueScrapeID: league_scrape_id,
			GameScrapeID: game_scrape_id,
			Group: group,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createLeagueSubdomain(info db.LeagueScrapeInfo) string {
	ret := fmt.Sprintf("/league/%d/%s/%d", info.SeasonStart, info.Type, info.ScrapeID)
	if info.Group != "" {
		ret += fmt.Sprintf("?group=%s", info.Group)
	}
	return ret
}

func getGameScrapeInfo(subdomain string) (int, string, int, int, string) {
	splits := strings.Split(subdomain, "/")

	season_start, _ := strconv.Atoi(splits[2])
	ltype := splits[3]
	league_scrape_id, _ := strconv.Atoi(splits[4])

	splits = strings.Split(splits[5], "?")
	game_scrape_id, _ := strconv.Atoi(splits[0])

	group := ""
	if len(splits) > 1 {
		splits = strings.Split(splits[1], "=")
		group = splits[1]
	}

	return season_start, ltype, league_scrape_id, game_scrape_id, group
}
