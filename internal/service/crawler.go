package service

import (
	"crawler/internal/scraper"
	"crawler/pkg/config"
	"crawler/pkg/cookies"
	"crawler/pkg/logger"
	"fmt"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
	"crawler/internal/repository"
)

type ICrawlerService interface {
	CheckPrerequisites() error
	ExecuteCrawl() error
	Initialize() error
	Cleanup()
}

type CrawlerService struct {
	config     *config.Config
	pw         *playwright.Playwright
	browser    playwright.Browser
	repository repository.ArticleRepository
}

func NewCrawlerService(cfg *config.Config, repo repository.ArticleRepository) ICrawlerService {
	return &CrawlerService{
		config:     cfg,
		repository: repo,
	}
}

// CheckPrerequisites 检查爬虫执行的前置条件
func (s *CrawlerService) CheckPrerequisites() error {
	// 检查 cookies 文件
	cookiesPath := s.config.App.CookiesFilePath
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		return fmt.Errorf("cookies文件不存在: %s", cookiesPath)
	}
	return nil
}

// Initialize 初始化浏览器
func (s *CrawlerService) Initialize() error {
	start := time.Now()
	logger.Info("初始化 Playwright",
		"step", "start",
		"timestamp", start.Format(time.RFC3339),
	)

	pw, err := playwright.Run()
	if err != nil {
		logger.Error("Playwright 安装失败",
			"error", err,
			"duration", time.Since(start).String(),
		)
		return fmt.Errorf("playwright installation failed: %w", err)
	}

	browserOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // 默认使用无头模式
	}

	browser, err := pw.Chromium.Launch(browserOpts)
	if err != nil {
		pw.Stop()
		logger.Error("浏览器启动失败",
			"error", err,
			"options", browserOpts,
			"duration", time.Since(start).String(),
		)
		return fmt.Errorf("browser launch failed: %w", err)
	}

	s.pw = pw
	s.browser = browser

	logger.Info("Playwright 初始化完成",
		"duration", time.Since(start).String(),
	)

	return nil
}

// Cleanup 清理资源
func (s *CrawlerService) Cleanup() {
	if s.browser != nil {
		s.browser.Close()
	}
	if s.pw != nil {
		s.pw.Stop()
	}
}

// ExecuteCrawl 执行爬虫任务
func (s *CrawlerService) ExecuteCrawl() error {
	start := time.Now()
	logger.Info("开始执行爬虫任务",
		"timestamp", start.Format(time.RFC3339),
	)

	// 初始化浏览器
	if err := s.Initialize(); err != nil {
		return err
	}
	defer s.Cleanup()

	// 创建新的上下文
	context, err := s.browser.NewContext()
	if err != nil {
		return fmt.Errorf("failed to create browser context: %w", err)
	}
	defer context.Close()

	// 加载 cookies
	if err := cookies.LoadCookies(context, s.config.App.CookiesFilePath); err != nil {
		logger.Warn("加载 Cookies 失败，将使用无登录模式",
			"error", err,
			"cookiesPath", s.config.App.CookiesFilePath,
		)
	}

	// 创建新页面
	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create new page: %w", err)
	}
	defer page.Close()

	// 访问目标页面
	targetURL := "https://www.zhihu.com/creator/manage/creation/article"
	logger.Info("开始访问目标页面",
		"url", targetURL,
	)

	if _, err := page.Goto(targetURL); err != nil {
		return fmt.Errorf("failed to navigate to %s: %w", targetURL, err)
	}

	// 提取数据
	data, err := scraper.ExtractData(page)
	if err != nil {
		return fmt.Errorf("failed to extract data: %w", err)
	}

	// 保存到数据库
	if err := s.repository.Save(data); err != nil {
		return fmt.Errorf("failed to save articles: %w", err)
	}

	logger.Info("数据提取完成",
		"articleCount", len(data),
		"duration", time.Since(start).String(),
	)

	return nil
}
