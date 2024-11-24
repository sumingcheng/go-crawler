package scraper

import (
	"crawler/pkg/logger"
	"log"
	"strconv"
	"strings"
	"time"

	playwright2 "github.com/playwright-community/playwright-go"
)

type ArticleCard struct {
	Title         string
	Link          string
	Description   string
	PublishedTime string
	Stats         ArticleStats
}

type ArticleStats struct {
	Reads     int
	Upvote    int
	Comments  int
	Bookmarks int
	Likes     int
}

func ExtractData(page playwright2.Page) ([]ArticleCard, error) {
	var articles []ArticleCard
	seenLinks := make(map[string]bool)
	noNewDataCount := 0 // 记录连续没有新数据的次数

	// 等待列表容器加载
	_, err := page.WaitForSelector("div[role='list']")
	if err != nil {
		return nil, err
	}

	logger.Info("开始提取文章数据")

	// 无限循环，直到确认没有新数据
	for {
		// 执行滚动
		_, err := page.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`)
		if err != nil {
			return nil, err
		}

		// 等待新内容加载（3秒）
		time.Sleep(3 * time.Second)

		// 获取当前所有文章卡片
		cards, err := page.QuerySelectorAll(".CreationManage-CreationCard")
		if err != nil {
			return nil, err
		}

		previousCount := len(articles)

		// 处理新的文章卡片
		for _, card := range cards {
			article, err := extractCardDetails(card)
			if err != nil {
				logger.Error("提取文章详情失败", "error", err)
				continue
			}

			// 检查是否已经处理过这篇文章
			if !seenLinks[article.Link] {
				seenLinks[article.Link] = true
				articles = append(articles, article)
				logger.Info("成功提取文章",
					"title", article.Title,
					"link", article.Link,
					"total_articles", len(articles),
				)
			}
		}

		// 检查是否有新数据
		if len(articles) == previousCount {
			noNewDataCount++
			logger.Info("本次滚动未发现新文章",
				"retry_count", noNewDataCount,
				"total_articles", len(articles),
			)
			// 连续3次没有新数据，认为已经到底
			if noNewDataCount >= 3 {
				logger.Info("已到达页面底部，停止提取",
					"total_articles", len(articles),
				)
				break
			}
		} else {
			// 有新数据，重置计数器
			noNewDataCount = 0
		}
	}

	return articles, nil
}

func extractCardDetails(card playwright2.ElementHandle) (ArticleCard, error) {
	var article ArticleCard
	var err error

	article.Title, err = getText(card, ".CreationCardTitle-wrapper")
	if err != nil {
		return article, err
	}

	linkElement, err := card.QuerySelector("a.css-959ia8")
	if err != nil {
		return article, err
	}
	article.Link, err = linkElement.GetAttribute("href")
	if err != nil {
		return article, err
	}

	article.Description, err = getText(card, ".CreationCardContent-text span")
	if err != nil {
		return article, err
	}

	article.PublishedTime, err = getText(card, ".css-zzavo4")
	if err != nil {
		return article, err
	}

	// Extract statistics
	stats, err := extractStats(card)
	if err != nil {
		return article, err
	}
	article.Stats = stats

	return article, nil
}

func getText(card playwright2.ElementHandle, selector string) (string, error) {
	element, err := card.QuerySelector(selector)
	if err != nil {
		return "", err
	}
	return element.InnerText()
}

func extractStats(card playwright2.ElementHandle) (ArticleStats, error) {
	var stats ArticleStats
	statElements, err := card.QuerySelectorAll(".css-150duks div")
	if err != nil {
		return stats, err
	}

	var lastNumber int
	for _, element := range statElements {
		textContent, err := element.InnerText()
		if err != nil {
			log.Printf("无法获取统计信息文本: %v", err)
			continue
		}
		textContent = strings.TrimSpace(textContent) // 清除两侧可能的空格

		if number, err := strconv.Atoi(textContent); err == nil {
			// 如果转换成功，说明是数字，记录下来
			lastNumber = number
		} else {
			// 如果转换失败，说明是文本标签，与上一个数字配对
			switch textContent {
			case "阅读":
				stats.Reads = lastNumber
			case "赞同":
				stats.Upvote = lastNumber
			case "评论":
				stats.Comments = lastNumber
			case "收藏":
				stats.Bookmarks = lastNumber
			case "喜欢":
				stats.Likes = lastNumber
			}
			lastNumber = 0 // 重置数字，避免错误关联
		}
	}
	return stats, nil
}
