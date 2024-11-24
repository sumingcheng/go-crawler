package main

import (
	"crawler/internal/controller"
	"crawler/internal/repository"
	"crawler/internal/router"
	"crawler/internal/service"
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"crawler/pkg/mysql"
	"log"
)

func main() {
	// 从配置文件加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化日志系统
	if err := logger.InitializeLogger(cfg.Logger); err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}

	// 初始化数据库连接
	db, err := mysql.NewDB(cfg.MySQL)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 初始化依赖
	repo := repository.NewMySQLArticleRepository(db)
	crawlerService := service.NewCrawlerService(cfg, repo)
	crawlerController := controller.NewCrawlerController(crawlerService)

	// 初始化路由
	r, err := router.NewRouter(cfg, crawlerController)
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
