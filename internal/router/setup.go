package router

import (
	"context"
	"crawler/internal/controller"
	"crawler/internal/middleware"
	"crawler/internal/service"
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	config  *config.Config
	engine  *gin.Engine
	crawler *controller.CrawlerController
}

// 配置检查
func validateConfig(cfg *config.Config) error {
	if cfg.Server.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout: %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout: %v", cfg.Server.WriteTimeout)
	}
	return nil
}

// 设置中间件
func setupMiddlewares(engine *gin.Engine, cfg *config.Config) {
	engine.Use(
		middleware.RequestID(),
		middleware.Cors(cfg),
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: customLogFormatter,
			SkipPaths: []string{"/health", "/metrics"},
		}),
		gin.Recovery(),
	)
}

// 自定义日志格式
func customLogFormatter(param gin.LogFormatterParams) string {
	// 使用结构化日志
	fields := map[string]interface{}{
		"status":     param.StatusCode,
		"method":     param.Method,
		"path":       param.Path,
		"client_ip":  param.ClientIP,
		"duration":   param.Latency.Seconds(),
		"user_agent": param.Request.UserAgent(),
		"error":      param.ErrorMessage,
		"request_id": param.Request.Header.Get("X-Request-ID"),
	}

	if param.StatusCode >= 400 {
		logger.Error("HTTP请求异常", fields)
	} else {
		logger.Info("HTTP请求", fields)
	}
	return ""
}

func NewRouter(cfg *config.Config) (*Router, error) {
	// 配置检查
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 设置 gin 模式
	gin.SetMode(cfg.Server.Mode)

	engine := gin.New()

	// 设置受信任代理
	if len(cfg.Server.TrustedProxies) > 0 {
		if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
			return nil, fmt.Errorf("set trusted proxies failed: %w", err)
		}
	}

	// 初始化服务和控制器
	crawlerService := service.NewCrawlerService(cfg)
	crawlerController := controller.NewCrawlerController(crawlerService)

	router := &Router{
		config:  cfg,
		engine:  engine,
		crawler: crawlerController,
	}

	// 设置中间件
	setupMiddlewares(engine, cfg)

	return router, nil
}

func (r *Router) Run(addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           r.engine,
		ReadTimeout:       r.config.Server.ReadTimeout,
		WriteTimeout:      r.config.Server.WriteTimeout,
		MaxHeaderBytes:    r.config.Server.MaxHeaderBytes,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("正在关闭服务器...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("服务器关闭异常", "error", err)
		}
	}()

	logger.Info("HTTP服务启动",
		"addr", addr,
		"mode", gin.Mode(),
		"read_timeout", r.config.Server.ReadTimeout,
		"write_timeout", r.config.Server.WriteTimeout,
		"max_header_bytes", r.config.Server.MaxHeaderBytes,
	)

	return srv.ListenAndServe()
}
