package router

import (
	"context"
	"crawler/internal/controller"
	"crawler/internal/middleware"
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
	config   *config.Config
	engine   *gin.Engine
	handlers *controller.Handlers
}

// NewRouter 创建并初始化 HTTP 路由实例
func NewRouter(cfg *config.Config, crawlerController *controller.CrawlerController) (*Router, error) {
	// 设置 gin 模式
	gin.SetMode(cfg.Server.Mode)

	ginEngine := gin.New()

	// 设置受信任代理服务器IP列表
	if len(cfg.Server.TrustedProxies) > 0 {
		if err := ginEngine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
			return nil, fmt.Errorf("设置受信任代理失败: %w", err)
		}
	}

	handlers := &controller.Handlers{
		Crawler: crawlerController,
	}

	ginRouter := &Router{
		config:   cfg,
		engine:   ginEngine,
		handlers: handlers,
	}

	// 全局中间件
	ginEngine.Use(
		middleware.TraceID(),
		middleware.Cors(cfg),
		gin.Logger(),
		gin.Recovery(),
	)

	return ginRouter, nil
}

// Run 启动 HTTP 服务器并支持优雅关闭
func (r *Router) Run(addr string) error {
	srv := &http.Server{
		Addr:           addr,
		Handler:        r.engine,
		ReadTimeout:    r.config.Server.ReadTimeout,
		WriteTimeout:   r.config.Server.WriteTimeout,
		MaxHeaderBytes: r.config.Server.MaxHeaderBytes,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("开始执行优雅关闭...")

		// 1. 停止接收新请求
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 2. 等待现有请求处理完成
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("HTTP服务关闭异常", "error", err)
		}

		logger.Info("服务已完全关闭")
	}()

	logger.Info("HTTP服务启动", "addr", addr, "mode", gin.Mode())
	return srv.ListenAndServe()
}

func (r *Router) SetupRoutes() {
	// 注册业务路由
	r.setupCrawlerRoutes()
	// 注册系统路由（健康检查等）
	r.setupHealthRoutes()
}
