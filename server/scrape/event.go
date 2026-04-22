package scrape

import (
	"database/sql"
	"fmt"
	"log"
	"polo-prophet/server/db"
	"strconv"
	"strings"
	// "time"

	"github.com/gocolly/colly/v2"
)

func ScrapeEvents(database *sql.DB, game_id int64) {
	c := colly.NewCollector()

	scrapeInfo_id, _ := db.GameScrapeInfoExists(database, game_id)
	if scrapeInfo_id == 0 {
		return
	}
	scrapeInfo, err := db.GetGameScrapeInfo(database, scrapeInfo_id)
	if err != nil {
		log.Fatal(err)
	}

	domain := "https://www.wasserball-team-deutschland.de" + createGameSubdomain(scrapeInfo)

	c.OnHTML(".timeline-item", func(h *colly.HTMLElement) {
		getEventInfo(h, database, int(game_id))
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			if exist, _ := db.EventsExists(database, int(game_id)); exist {
				db.DeleteEventsForGame(database, game_id)
			}
		}
	})

	err = c.Visit(domain)
	if err != nil {
		log.Print(err)
	}
}

func getEventInfo(e *colly.HTMLElement, database *sql.DB, game_id int) {
	team := strings.TrimPrefix(e.Attr("class"), "timeline-item ")

	player := e.ChildText("h3")

	etype := e.ChildText(".timeline-top span.muted")

	pt := strings.Split(e.ChildText(".timeline-top span.pill"), " · ")
	var period int
	var time int
	var err error
	if len(pt) >= 2 {
		period, err = strconv.Atoi(strings.TrimPrefix(pt[0], "Abschnitt "))
		if err != nil {
			log.Fatal(err)
		}
		time_segments := strings.Split(pt[1], ":")
		if len(time_segments) >= 2 {
			seconds, err := strconv.Atoi(time_segments[1])
			if err != nil {
				log.Fatal(err)
			}
			minutes, err := strconv.Atoi(time_segments[0])
			if err != nil {
				log.Fatal(err)
			}
			time = minutes * 60 + seconds
		} else {
			time = 0
		}
	} else {
		period = 0
		time = 0
	}

	db.CreateEvent(database, db.Event{
		ID:     0,
		GameID: game_id,
		Team:   team,
		Player: player,
		Type:   etype,
		Period: period,
		Time:   time,
	})
}

func createGameSubdomain(info db.GameScrapeInfo) string {
	ret := fmt.Sprintf("/game/%d/%s/%d/%d", info.SeasonStart, info.Type, info.LeagueScrapeID, info.GameScrapeID)
	if info.Group != "" {
		ret += fmt.Sprintf("?group=%s", info.Group)
	}
	return ret
}
