// playwright.go
package main

import (
	"crawler/common/playwright"
	"crawler/config"
	"log"
)

func main() {
	log.Println("开始执行程序")

	// 从配置文件加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 选择使用Playwright爬取数据
	playwright.ExecutePlaywright(cfg)

	// 未来添加HTTP请求方式
	// if cfg.UseHTTP {
	//     executeHTTP(cfg)
	// }
}
