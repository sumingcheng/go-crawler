package service

import (
	"crawler/internal/playwright"
	"crawler/pkg/config"
	"fmt"
	"os"
)

type CrawlerService struct {
	config *config.Config
}

func NewCrawlerService(cfg *config.Config) *CrawlerService {
	return &CrawlerService{
		config: cfg,
	}
}

// CheckPrerequisites 检查爬虫执行的前置条件
func (s *CrawlerService) CheckPrerequisites() error {
	// 检查 cookies 文件
	cookiesPath := s.config.App.CookiesFilePath
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		return fmt.Errorf("cookies文件不存在: %s", cookiesPath)
	}
	return nil
}

func (s *CrawlerService) ExecuteCrawl() error {
	return playwright.ExecutePlaywright(s.config)
}
