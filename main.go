package main

import (
	"fmt"
	"log"
	
	"property-scraper/scraper"
)

func main() {
	fmt.Println("Starting property scraper...")
	
	err := scraper.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Scraping completed done")
}
