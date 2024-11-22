package service

import (
	"crawler/internal/playwright"
	"crawler/pkg/config"
)

type CrawlerService struct {
	config *config.Config
}

func NewCrawlerService(cfg *config.Config) *CrawlerService {
	return &CrawlerService{
		config: cfg,
	}
}

func (s *CrawlerService) ExecuteCrawl() error {
	return playwright.ExecutePlaywright(s.config)
}
