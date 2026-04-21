package main

import (
	"log"
	"polo-prophet/server/api"
	"polo-prophet/server/db"
	"polo-prophet/server/scrape"
)

func main() {
	database := db.Connect()
	server := api.NewServer(database)

	go server.Start(":8080")

	go func ()  {
		scrape.ScrapeLeagues(database)

		league_ids, err := db.GetAllLeagueIds(database)
		if err != nil {
			log.Fatal(err)
		}
		for _, id := range league_ids {
			scrape.ScrapeGames(database, id)
		}
	}()

	select {}
}
