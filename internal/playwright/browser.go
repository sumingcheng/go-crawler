package playwright

import (
	"crawler/internal/scraper"
	"crawler/pkg/config"
	"crawler/pkg/cookies"
	"crawler/pkg/logger"
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
)

func InitializePlaywright() (*playwright.Playwright, playwright.Browser, error) {
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
		return nil, nil, fmt.Errorf("playwright installation failed: %w", err)
	}

	browserOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	}

	browser, err := pw.Chromium.Launch(browserOpts)
	if err != nil {
		pw.Stop()
		logger.Error("浏览器启动失败",
			"error", err,
			"options", browserOpts,
			"duration", time.Since(start).String(),
		)
		return nil, nil, fmt.Errorf("browser launch failed: %w", err)
	}

	logger.Info("Playwright 初始化完成",
		"duration", time.Since(start).String(),
	)

	return pw, browser, nil
}

func ExecutePlaywright(cfg *config.Config) error {
	start := time.Now()
	logger.Info("开始执行爬虫任务",
		"timestamp", start.Format(time.RFC3339),
	)

	pw, browser, err := InitializePlaywright()
	if err != nil {
		return fmt.Errorf("playwright initialization failed: %w", err)
	}

	defer func() {
		browser.Close()
		pw.Stop()
		logger.Info("爬虫任务结束",
			"duration", time.Since(start).String(),
		)
	}()

	context, err := browser.NewContext()
	if err != nil {
		return fmt.Errorf("failed to create browser context: %w", err)
	}
	defer context.Close()

	if err := cookies.LoadCookies(context, cfg.App.CookiesFilePath); err != nil {
		logger.Warn("加载 Cookies 失败，将使用无登录模式",
			"error", err,
			"cookiesPath", cfg.App.CookiesFilePath,
		)
	}

	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("failed to create new page: %w", err)
	}
	defer page.Close()

	targetURL := "https://www.zhihu.com/creator/manage/creation/article"
	logger.Info("开始访问目标页面",
		"url", targetURL,
	)

	if _, err := page.Goto(targetURL); err != nil {
		return fmt.Errorf("failed to navigate to %s: %w", targetURL, err)
	}

	data, err := scraper.ExtractData(page)
	if err != nil {
		return fmt.Errorf("failed to extract data: %w", err)
	}

	logger.Info("数据提取完成",
		"articleCount", len(data),
		"duration", time.Since(start).String(),
	)

	return nil
}
