package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 注册所有路由
func (r *Router) SetupRoutes() {
	// 注册健康检查路由
	r.setupHealthRoutes()

	// 注册 API 路由
	r.setupAPIRoutes()
}

// setupHealthRoutes 健康检查路由
func (r *Router) setupHealthRoutes() {
	r.engine.GET("/health", r.handleHealth)
}

// setupAPIRoutes API相关路由
func (r *Router) setupAPIRoutes() {
	api := r.engine.Group("/api")
	{
		// 保持原有的路由
		api.POST("/crawl", r.crawler.HandleCrawl)

		// 后续可以添加更多路由组
		// r.setupUserRoutes(api)
		// r.setupDataRoutes(api)
	}
}

// handleHealth 健康检查处理器
func (r *Router) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "up",
		"time":   time.Now().Format(time.RFC3339),
	})
}
