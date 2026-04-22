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

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("GameID %d, StatusCode %d\n", game_id, r.StatusCode)
		if r.StatusCode == 429 {
			retryAfter := r.Headers.Get("Retry-After")
			wait := 5 * time.Second
			if secs, err := strconv.Atoi(retryAfter); err == nil {
				wait = time.Duration(secs) * time.Second
			}
			log.Printf("EventScraping for Game_ID %d: Retrying after %dseconds ...\n", game_id, wait/1000000000)
			time.Sleep(wait)
			r.Request.Retry()
		}
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("GameID %d, StatusCode %d\n", game_id, r.StatusCode)
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

	db.CreateEvent(database, db.Event{
		ID:     0,
		GameID: game_id,
		Team:   team,
		Player: player,
		Type:   etype,
	})
}

func createGameSubdomain(info db.GameScrapeInfo) string {
	ret := fmt.Sprintf("/game/%d/%s/%d/%d", info.SeasonStart, info.Type, info.LeagueScrapeID, info.GameScrapeID)
	if info.Group != "" {
		ret += fmt.Sprintf("?group=%s", info.Group)
	}
	return ret
}
