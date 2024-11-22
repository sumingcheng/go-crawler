package router

import (
	"crawler/internal/playwright"
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	config *config.Config
	engine *gin.Engine
}

func NewRouter(cfg *config.Config) *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{
		config: cfg,
		engine: engine,
	}
}

func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api")
	{
		api.POST("/crawl", r.handleCrawl)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r *Router) handleCrawl(c *gin.Context) {
	logger.Info("收到爬取请求")

	if err := playwright.ExecutePlaywright(r.config); err != nil {
		logger.Error("爬取失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "爬取成功",
	})
}
