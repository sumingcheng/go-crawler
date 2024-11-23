package controller

import (
	"crawler/internal/service"
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"crawler/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CrawlerController struct {
	crawlerService *service.CrawlerService
}

func NewCrawlerHandler(cfg *config.Config) *CrawlerController {
	return &CrawlerController{
		crawlerService: service.NewCrawlerService(cfg),
	}
}

func (cc *CrawlerController) HandleCrawl(c *gin.Context) {
	start := time.Now()
	logger.Info("收到爬取请求",
		"trace_id", c.GetString("trace_id"),
	)

	if err := cc.crawlerService.CheckPrerequisites(); err != nil {
		logger.Error("前置条件检查失败",
			"error", err,
			"duration", time.Since(start).String(),
			"trace_id", c.GetString("trace_id"),
		)
		response.Error(c, http.StatusBadRequest, "爬取前置条件不满足: "+err.Error())
		return
	}

	if err := cc.crawlerService.ExecuteCrawl(); err != nil {
		logger.Error("爬取失败",
			"error", err,
			"duration", time.Since(start).String(),
			"trace_id", c.GetString("trace_id"),
		)
		response.Error(c, http.StatusInternalServerError, "爬取失败: "+err.Error())
		return
	}

	logger.Info("爬取完成",
		"duration", time.Since(start).String(),
		"trace_id", c.GetString("trace_id"),
	)

	response.Success(c, "爬取成功", nil)
}
