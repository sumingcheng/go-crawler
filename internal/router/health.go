package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

// setupHealthRoutes 注册健康检查相关路由
func (r *Router) setupHealthRoutes() {
	r.engine.GET("/health", handleHealth)
}

// handleHealth 健康检查处理器
func handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "up",
		"time":   time.Now().Format(time.RFC3339),
	})
}
