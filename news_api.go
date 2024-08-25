package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gocolly/colly"
)

type Hackernews struct {
	Status       string     `json:"status"`
	TotalResults int        `json:"totalResults"`
	Articles     []Articles `json:"articles"`
}
type Source struct {
	ID   any    `json:"id"`
	Name string `json:"name"`
}
type Articles struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type NewsItem struct {
	Title   string
	Content string
	URL     string
}

func get_news(source string) ([]NewsItem, error) {
	date := time.Now()
	api := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&from=%s&sortBy=publishedAt&language=en&apiKey=e99372e12b6a4933a13b10e9ea2a7f9d", source, date)

	news_response, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer news_response.Body.Close() // Ensure the response body is closed

	if news_response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Response: %d", news_response.StatusCode)
	}

	bodyByte, err := io.ReadAll(news_response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading API Body: %s", err)
	}

	var hackernews Hackernews
	err = json.Unmarshal(bodyByte, &hackernews)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal API Body: %s", err)
	}

	var newsItems []NewsItem
	for i, article := range hackernews.Articles {
		if i < 60 {
			content, err := scraper(article.URL)
			if err != nil {
				return nil, err
			}
			newsItems = append(newsItems, NewsItem{
				Title:   article.Title,
				Content: content,
				URL:     article.URL,
			})
		} else {
			break
		}
	}

	return newsItems, nil
}

func scraper(url string) (string, error) {
	c := colly.NewCollector()
	var paragraphs []string
	c.OnHTML("article", func(e *colly.HTMLElement) {
		paragraphs = append(paragraphs, e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", paragraphs), nil
}
