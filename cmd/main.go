// playwright.go
package main

import (
	"crawler/common/playwright"
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

	// 使用封装的日志记录
	logger.Info("程序启动")

	// 记录带字段的日志
	logger.WithFields(map[string]interface{}{
		"username":     cfg.App.Username,
		"cookies_path": cfg.App.CookiesFilePath,
	}).Info("配置加载成功")

	// 选择使用Playwright爬取数据
	playwright.ExecutePlaywright(cfg)

	logger.Info("程序退出")
}
