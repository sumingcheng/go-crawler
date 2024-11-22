package controller

import (
	"crawler/internal/service"
	"crawler/pkg/logger"
	"crawler/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CrawlerController struct {
	crawlerService *service.CrawlerService
}

func NewCrawlerController(crawlerService *service.CrawlerService) *CrawlerController {
	return &CrawlerController{
		crawlerService: crawlerService,
	}
}

func (cc *CrawlerController) HandleCrawl(c *gin.Context) {
	start := time.Now()
	logger.Info("收到爬取请求",
		"trace_id", c.GetString("trace_id"),
	)

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
