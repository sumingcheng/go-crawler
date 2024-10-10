// main.go
package main

import (
	"crawler/config"
	"crawler/cookies"
	"crawler/playwright"
	playwright2 "github.com/playwright-community/playwright-go"
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

	extractData(err, page)

	log.Println("数据提取完成")
}

func extractData(err error, page playwright2.Page) {
	// 等待页面加载相关元素
	_, err = page.WaitForSelector("div[role='list']")
	if err != nil {
		log.Fatalf("等待列表元素加载失败: %v", err)
	}

	// 提取所有卡片信息
	cards, err := page.QuerySelectorAll(".CreationManage-CreationCard")
	if err != nil {
		log.Fatalf("查询文章卡片失败: %v", err)
	}

	for _, card := range cards {
		titleElement, err := card.QuerySelector(".CreationCardTitle-wrapper")
		if err != nil {
			log.Printf("查询标题元素失败: %v", err)
			continue
		}
		title, err := titleElement.InnerText()
		if err != nil {
			log.Printf("获取标题文本失败: %v", err)
			continue
		}

		linkElement, err := card.QuerySelector("a.css-959ia8")
		if err != nil {
			log.Printf("查询链接元素失败: %v", err)
			continue
		}
		link, err := linkElement.GetAttribute("href")
		if err != nil {
			log.Printf("获取链接属性失败: %v", err)
			continue
		}

		descriptionElement, err := card.QuerySelector(".CreationCardContent-text span")
		if err != nil {
			log.Printf("查询描述元素失败: %v", err)
			continue
		}
		description, err := descriptionElement.InnerText()
		if err != nil {
			log.Printf("获取描述文本失败: %v", err)
			continue
		}

		publishedElement, err := card.QuerySelector(".css-zzavo4")
		if err != nil {
			log.Printf("查询发布时间元素失败: %v", err)
			continue
		}
		publishedTime, err := publishedElement.InnerText()
		if err != nil {
			log.Printf("获取发布时间文本失败: %v", err)
			continue
		}

		log.Printf("标题: %s\n链接: %s\n描述: %s\n发布时间: %s\n", title, link, description, publishedTime)
	}
}
