package main

import (
	"crawler/internal/di"
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"crawler/pkg/mysql"
	"log"
	"os"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 2. 初始化日志系统
	if err := logger.InitializeLogger(cfg.Logger); err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 3. 初始化数据库连接
	db, err := mysql.NewDB(cfg.MySQL)
	if err != nil {
		log.Fatal("数据库连接失败", "error", err)
	}

	// 4. 初始化依赖注入容器
	container, err := di.NewContainer(cfg, db)
	if err != nil {
		log.Fatal("依赖注入容器初始化失败", "error", err)
	}

	// 确保资源正确清理
	defer container.Cleanup()

	// 5. 设置路由
	container.Router.SetupRoutes()

	// 6. 启动服务
	logger.Info("开始启动服务", "port", cfg.Server.Port)
	if err := container.Router.Run(cfg.Server.Port); err != nil {
		log.Fatal("服务启动失败", "error", err)
		os.Exit(1)
	}
}
