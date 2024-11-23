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
func NewRouter(cfg *config.Config) (*Router, error) {
	// 设置 gin 模式
	gin.SetMode(cfg.Server.Mode)

	engine := gin.New()

	// 设置受信任代理服务器IP列表
	if len(cfg.Server.TrustedProxies) > 0 {
		if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
			return nil, fmt.Errorf("设置受信任代理失败: %w", err)
		}
	}

	router := &Router{
		config:   cfg,
		engine:   engine,
		handlers: controller.NewHandlers(cfg),
	}

	// 设置中间件
	setupGlobalMiddlewares(engine, cfg)

	return router, nil
}

// setupGlobalMiddlewares 配置全局中间件
func setupGlobalMiddlewares(engine *gin.Engine, cfg *config.Config) {
	engine.Use(
		middleware.TraceID(),
		middleware.Cors(cfg),
		gin.Logger(),
		gin.Recovery(),
	)
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
	)

	return srv.ListenAndServe()
}
