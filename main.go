package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
	"time"
)

// Movie holds the movie object where the scraped films will be stored
type Movie struct {
	Title   string
	Details string
	Theatre string
	Date    time.Time
}

func scrape() []Movie {

	daysTranslate := strings.NewReplacer(
		"Lunes", "Monday",
		"Martes", "Tuesday",
		"Miercoles", "Wednesday",
		"Miércoles", "Wednesday",
		"Jueves", "Thursday",
		"Viernes", "Friday",
		"Sabado", "Saturday",
		"Sábado", "Saturday",
		"Domingo", "Sunday",
	)

	movies := []Movie{}
	c := colly.NewCollector()
	detailCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(".event-titulo", func(e *colly.HTMLElement) {
		events := e.ChildAttrs("a", "href")
		event := events[len(events)-1]
		detailCollector.Visit("http://www.phenomena-experience.com/" + event)
	})

	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {

		if e.DOM.Find(".pase-hora").Length() > 0 {

			dateStr := daysTranslate.Replace(e.DOM.Find(".pase-fecha").First().Text())
			dateTime := strings.Fields(e.DOM.Find(".pase-hora").First().Text())[0]
			movie := Movie{}
			movie.Title = e.DOM.Find(".titulo").Find("span").Text()
			movie.Details = e.ChildText(".datos2")
			movie.Theatre = "Phenomena Experience"

			if t, err := time.Parse("Monday, 02/01/2006 15:04h", fmt.Sprintf("%s %s", dateStr, dateTime)); err != nil {
				log.Printf("Error: %s\n", err)
			} else {
				movie.Date = t.Add(-time.Hour * 2)
			}
			movies = append(movies, movie)
		}
	})
	c.OnScraped(func(r *colly.Response) {
		log.Println("Scrape finished")
	})

	c.Visit("http://www.phenomena-experience.com/programacion-mensual/todo.html")

	return movies
}

func main() {
	log.Println(scrape())
}
