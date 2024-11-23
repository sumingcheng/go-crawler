package api

import (
	"crawler/internal/controller"
	"crawler/internal/service"
	"crawler/pkg/config"
)

// Handlers 集中管理所有 API 处理器
type Handlers struct {
	Crawler *controller.CrawlerController
	// 可以添加更多控制器
}

// NewHandlers 创建 API 处理器集合
func NewHandlers(cfg *config.Config) *Handlers {
	// 初始化服务
	crawlerService := service.NewCrawlerService(cfg)

	return &Handlers{
		Crawler: controller.NewCrawlerController(crawlerService),
	}
}
