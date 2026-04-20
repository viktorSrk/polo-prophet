package main

import (
	"polo-prophet/server/api"
	"polo-prophet/server/db"
	"polo-prophet/server/scrape"
)

func main() {
	database := db.Connect()
	server := api.NewServer(database)

	go server.Start(":8080")

	go scrape.ScrapeLeagues(database)

	select {}
}
