package playwright

import (
	"crawler/common/scraper"
	"crawler/config"
	"crawler/cookies"
	"fmt"
	"log"
)

// ExecutePlaywright 封装了使用Playwright进行网页爬取的逻辑
func ExecutePlaywright(cfg config.Config) {
	pw, browser, err := InitializePlaywright()
	if err != nil {
		log.Fatalf("Playwright 初始化失败: %v", err)
	}
	defer browser.Close()
	defer pw.Stop()
	log.Println("Playwright 和浏览器初始化成功")

	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("创建浏览器上下文失败: %v", err)
	}
	defer context.Close()
	log.Println("浏览器上下文创建成功")

	if err := cookies.LoadCookies(context, cfg.CookiesFilePath); err != nil {
		log.Fatalf("加载Cookies失败: %v", err)
	}
	log.Println("Cookies 加载成功")

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("创建页面失败: %v", err)
	}
	defer page.Close()
	log.Println("页面创建成功")

	if _, err := page.Goto("https://www.zhihu.com/creator/manage/creation/article"); err != nil {
		log.Fatalf("导航到页面失败: %v", err)
	} else {
		log.Println("成功导航到指定页面")
	}

	data, err := scraper.ExtractData(page)
	if err != nil {
		log.Fatalf("数据提取失败: %v", err)
	}
	fmt.Println(data)
	log.Println("数据提取完成")
}
