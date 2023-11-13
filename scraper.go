package main

import (
	"log"
	"time"

	"github.com/brus-fabrika/rssapi/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	fetchInterval time.Duration,
) {

	log.Printf("Starting scraper with %v workers with %s interval", concurrency, fetchInterval)
	ticker := time.NewTicker(fetchInterval)
	for ; ; <-ticker.C {
		log.Println("Fetching feeds")
	}
}
