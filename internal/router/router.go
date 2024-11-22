package router

import (
	"crawler/internal/controller"
	"crawler/internal/middleware"
	"crawler/internal/service"
	"crawler/pkg/config"
	"crawler/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Router struct {
	config  *config.Config
	engine  *gin.Engine
	crawler *controller.CrawlerController
}

func NewRouter(cfg *config.Config) *Router {
	// 设置 gin 的日志输出
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	// 使用自定义日志格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Info("路由注册",
			"method", httpMethod,
			"path", absolutePath,
			"handler", handlerName,
			"middleware_count", nuHandlers,
		)
	}

	engine := gin.New()

	// 初始化服务
	crawlerService := service.NewCrawlerService(cfg)
	crawlerController := controller.NewCrawlerController(crawlerService)

	// 添加中间件
	engine.Use(
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			// 记录到我们的日志系统
			logger.Info("HTTP请求",
				"status", param.StatusCode,
				"method", param.Method,
				"path", param.Path,
				"client_ip", param.ClientIP,
				"duration", param.Latency,
				"user_agent", param.Request.UserAgent(),
				"error", param.ErrorMessage,
			)
			return "" // 不输出到控制台
		}),
		gin.Recovery(),
		middleware.Trace(),
		middleware.Cors(),
	)

	return &Router{
		config:  cfg,
		engine:  engine,
		crawler: crawlerController,
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Router) Run(addr string) error {
	logger.Info("HTTP服务启动",
		"addr", addr,
		"mode", gin.Mode(),
	)
	return r.engine.Run(addr)
}
