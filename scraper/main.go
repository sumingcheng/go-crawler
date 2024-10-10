package scraper

import (
	playwright2 "github.com/playwright-community/playwright-go"
	"log"
)

type ArticleCard struct {
	Title         string
	Link          string
	Description   string
	PublishedTime string
}

func ExtractData(page playwright2.Page) ([]ArticleCard, error) {
	var articles []ArticleCard

	_, err := page.WaitForSelector("div[role='list']")
	if err != nil {
		return nil, err
	}

	cards, err := page.QuerySelectorAll(".CreationManage-CreationCard")
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		article, err := extractCardDetails(card)
		if err != nil {
			log.Printf("Error extracting card details: %v", err)
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func extractCardDetails(card playwright2.ElementHandle) (ArticleCard, error) {
	var article ArticleCard
	var err error

	titleElement, err := card.QuerySelector(".CreationCardTitle-wrapper")
	if err != nil {
		return article, err
	}
	article.Title, err = titleElement.InnerText()
	if err != nil {
		return article, err
	}

	linkElement, err := card.QuerySelector("a.css-959ia8")
	if err != nil {
		return article, err
	}
	article.Link, err = linkElement.GetAttribute("href")
	if err != nil {
		return article, err
	}

	descriptionElement, err := card.QuerySelector(".CreationCardContent-text span")
	if err != nil {
		return article, err
	}
	article.Description, err = descriptionElement.InnerText()
	if err != nil {
		return article, err
	}

	publishedElement, err := card.QuerySelector(".css-zzavo4")
	if err != nil {
		return article, err
	}
	article.PublishedTime, err = publishedElement.InnerText()
	if err != nil {
		return article, err
	}

	return article, nil
}
