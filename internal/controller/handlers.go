package controller

import "crawler/pkg/config"

// Handlers 集中管理所有控制器
type Handlers struct {
	Crawler *CrawlerController
	// 可以添加更多控制器
}

// NewHandlers 创建并初始化所有控制器
func NewHandlers(cfg *config.Config) *Handlers {
	return &Handlers{
		Crawler: NewCrawlerHandler(cfg),
	}
}
