package controller

import (
	"crawler/internal/service"
	"crawler/pkg/config"
)

// Handlers 集中管理所有控制器
type Handlers struct {
	Crawler *CrawlerController
	// 可以添加更多控制器
}

func InitializeHandlers(cfg *config.Config, crawlerService service.ICrawlerService) *Handlers {
	return &Handlers{
		Crawler: NewCrawlerController(crawlerService),
	}
}
