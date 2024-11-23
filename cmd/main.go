package main

import (
	"crawler/internal/router"
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"log"
)

func main() {
	// 从配置文件加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化日志系统
	if err := logger.Init(cfg.Logger); err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 记录关键配置信息
	logger.WithFields(map[string]interface{}{
		"username":     cfg.App.Username,
		"cookies_path": cfg.App.CookiesFilePath,
		"server_port":  cfg.Server.Port,
	}).Info("系统初始化完成")

	// 初始化路由
	r, err := router.NewRouter(cfg)
	if err != nil {
		logger.Error("路由初始化失败", "error", err)
		return
	}
	r.SetupRoutes()

	// 启动服务器
	if err := r.Run(cfg.Server.Port); err != nil {
		logger.Error("服务器启动失败", "error", err)
	}
}
