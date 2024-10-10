// main.go
package main

import (
	"crawler/config"
	"crawler/cookies"
	"crawler/playwright"
	"crawler/scraper"
	"fmt"
	"log"
)

func main() {
	log.Println("开始执行程序")

	// 从配置文件加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化Playwright并启动浏览器
	pw, browser, err := playwright.InitializePlaywright()
	if err != nil {
		log.Fatalf("Playwright 初始化失败: %v", err)
	}
	defer browser.Close()
	defer pw.Stop()
	log.Println("Playwright 和浏览器初始化成功")

	// 创建浏览器上下文
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("创建浏览器上下文失败: %v", err)
	}
	defer context.Close()
	log.Println("浏览器上下文创建成功")

	// 加载Cookies
	if err := cookies.LoadCookies(context, cfg.CookiesFilePath); err != nil {
		log.Fatalf("加载Cookies失败: %v", err)
	}
	log.Println("Cookies 加载成功")

	// 创建新的页面
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("创建页面失败: %v", err)
	}
	defer page.Close()
	log.Println("页面创建成功")

	// 导航到特定页面
	if _, err := page.Goto("https://www.zhihu.com/creator/manage/creation/article"); err != nil {
		log.Fatalf("导航到页面失败: %v", err)
	} else {
		log.Println("成功导航到指定页面")
	}

	// 提取数据
	data, err := scraper.ExtractData(page)
	if err != nil {
		return
	}
	fmt.Println(data)
	log.Println("数据提取完成")
}
