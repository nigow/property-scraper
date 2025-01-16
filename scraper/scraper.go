package scraper

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/PuerkitoBio/goquery"
)

func getEnvVariable(key string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load .env file: %w", err)
	}
	
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable not set")
	}
	return value, nil
}

type Property struct {
	Title    string
	Price    string
	Address  string
	Area     string
	Layout   string
	Age      string
	Station  string
	WalkTime string
	URL      string
}

func Run() error {
	// Create CSV file
	file, err := os.Create("properties.csv")
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV header
	writer.Write([]string{"Title", "Price", "Address", "Area", "Layout", "Age", "Station", "Walk Time", "URL"})

	baseURL, err := getEnvVariable("BASEURL")
	if err != nil {
		return fmt.Errorf("failed to get base URL: %w", err)
	}

	// Start scraping
	page := 1
	for {
		fullURL, err := getEnvVariable("URL")
		if err != nil {
			return fmt.Errorf("failed to get URL: %w", err)
		}
		url := fullURL
		if page > 1 {
			url += fmt.Sprintf("&page=%d", page)
		}

		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch page %d: %w", page, err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return fmt.Errorf("failed to parse page %d: %w", page, err)
		}

		if doc.Find(".cassetteitem").Length() == 0 {
			break
		}

		doc.Find(".cassetteitem").Each(func(i int, s *goquery.Selection) {
			title := s.Find(".cassetteitem_content-title").Text()
			price := s.Find(".cassetteitem_price--rent").Text()
			address := s.Find(".cassetteitem_detail-col1").Text()
			area := s.Find(".cassetteitem_madori").Text()
			layout := s.Find(".cassetteitem_madori").Text()
			age := s.Find(".cassetteitem_detail-col2").Text()
			station := s.Find(".cassetteitem_detail-col2").Text()
			walkTime := s.Find(".cassetteitem_detail-col2").Text()
			url, _ := s.Find(".cassetteitem_other-linktext").Attr("href")

			title = strings.TrimSpace(title)
			price = strings.TrimSpace(price)
			address = strings.TrimSpace(address)
			area = strings.TrimSpace(area)
			layout = strings.TrimSpace(layout)
			age = strings.TrimSpace(age)
			station = strings.TrimSpace(station)
			walkTime = strings.TrimSpace(walkTime)
			url = baseURL + strings.TrimSpace(url)

			writer.Write([]string{
				title,
				price,
				address,
				area,
				layout,
				age,
				station,
				walkTime,
				url,
			})
		})

		// pagination
		if doc.Find(".pagination-parts a:contains('次へ')").Length() == 0 {
			break
		}

		page++
		time.Sleep(1 * time.Second) 
	}

	return nil
}
