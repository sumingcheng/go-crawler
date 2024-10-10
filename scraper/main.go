package scraper

import (
	"log"

	playwright2 "github.com/playwright-community/playwright-go"
)

func ExtractData(page playwright2.Page) {
	// 等待页面加载相关元素
	_, err := page.WaitForSelector("div[role='list']")
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
