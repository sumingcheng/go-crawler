package controller

import (
	"github.com/gin-gonic/gin"
)

// ICrawlerController 爬虫控制器接口
type ICrawlerController interface {
	HandleCrawl(c *gin.Context)
}

// IHandlers 处理器集合接口
type IHandlers interface {
	GetCrawlerController() ICrawlerController
}

// Handlers 集中管理所有控制器
type Handlers struct {
	Crawler ICrawlerController
}

// GetCrawlerController 获取爬虫控制器
func (h *Handlers) GetCrawlerController() ICrawlerController {
	return h.Crawler
}

// NewHandlers 创建处理器集合
func NewHandlers(crawlerController ICrawlerController) IHandlers {
	return &Handlers{
		Crawler: crawlerController,
	}
}
