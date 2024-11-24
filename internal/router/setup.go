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
	config     *config.Config
	engine     *gin.Engine
	controller controller.ICrawlerController
}

func NewRouter(cfg *config.Config, controller controller.ICrawlerController) (*Router, error) {
	gin.SetMode(cfg.Server.Mode)

	ginEngine := gin.New()

	if err := setupMiddleware(cfg, ginEngine); err != nil {
		return nil, err
	}

	router := &Router{
		config:     cfg,
		engine:     ginEngine,
		controller: controller,
	}

	router.setupCrawlerRoutes()
	router.setupHealthRoutes()

	return router, nil
}

func setupMiddleware(cfg *config.Config, engine *gin.Engine) error {
	if len(cfg.Server.TrustedProxies) > 0 {
		if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
			return fmt.Errorf("设置受信任代理失败: %w", err)
		}
	}

	engine.Use(
		middleware.TraceID(),
		middleware.Cors(cfg),
		gin.Logger(),
		gin.Recovery(),
	)

	return nil
}

func (r *Router) ServeHTTP(addr string) error {
	srv := &http.Server{
		Addr:           addr,
		Handler:        r.engine,
		ReadTimeout:    r.config.Server.ReadTimeout,
		WriteTimeout:   r.config.Server.WriteTimeout,
		MaxHeaderBytes: r.config.Server.MaxHeaderBytes,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("开始执行优雅关闭...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("HTTP服务关闭异常", "error", err)
		}

		logger.Info("服务已完全关闭")
	}()

	logger.Info("HTTP服务启动", "addr", addr, "mode", gin.Mode())
	return srv.ListenAndServe()
}
