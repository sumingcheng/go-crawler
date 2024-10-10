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
	// 合并 defer 语句，统一处理关闭浏览器和停止 Playwright 的操作
	defer func() {
		browser.Close()
		pw.Stop()
	}()

	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("创建浏览器上下文失败: %v", err)
	}
	defer context.Close()

	if err := cookies.LoadCookies(context, cfg.CookiesFilePath); err != nil {
		log.Fatalf("加载Cookies失败: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("创建页面失败: %v", err)
	}
	defer page.Close()

	if _, err := page.Goto("https://www.zhihu.com/creator/manage/creation/article"); err != nil {
		log.Fatalf("导航到页面失败: %v", err)
	}

	data, err := scraper.ExtractData(page)
	if err != nil {
		log.Fatalf("数据提取失败: %v", err)
	}
	fmt.Println(data)
	log.Println("数据提取完成")
}
