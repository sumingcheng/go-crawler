package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupRoutes() {
	// 注册业务路由
	r.setupCrawlerRoutes()
	// 注册系统路由（健康检查等）
	r.setupHealthRoutes()
}

// setupHealthRoutes 注册健康检查相关路由
func (r *Router) setupHealthRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

// setupCrawlerRoutes 爬虫相关路由
func (r *Router) setupCrawlerRoutes() {
	api := r.engine.Group("/api")
	crawler := api.Group("/crawler")
	{
		crawler.POST("/zhihu", r.handlers.GetCrawlerController().HandleCrawl)
	}
}
