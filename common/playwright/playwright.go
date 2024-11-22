package playwright

import (
	"crawler/common/scraper"
	"crawler/cookies"
	"crawler/pkg/config"
	"crawler/pkg/logger"

	"github.com/playwright-community/playwright-go"
)

func InitializePlaywright() (*playwright.Playwright, playwright.Browser, error) {
	logger.Info("开始安装 Playwright")
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Playwright 启动成功")

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		pw.Stop()
		return nil, nil, err
	}
	logger.Info("浏览器启动成功")

	return pw, browser, nil
}

// ExecutePlaywright 封装了使用Playwright进行网页爬取的逻辑
func ExecutePlaywright(cfg *config.Config) error {
	pw, browser, err := InitializePlaywright()
	if err != nil {
		logger.Error("Playwright 初始化失败", "error", err)
		return err
	}

	defer func() {
		browser.Close()
		pw.Stop()
	}()

	context, err := browser.NewContext()
	if err != nil {
		logger.Error("创建浏览器上下文失败", "error", err)
		return err
	}
	defer context.Close()

	if err := cookies.LoadCookies(context, cfg.App.CookiesFilePath); err != nil {
		logger.Error("加载Cookies失败", "error", err)
		return err
	}

	page, err := context.NewPage()
	if err != nil {
		logger.Error("创建页面失败", "error", err)
		return err
	}
	defer page.Close()

	if _, err := page.Goto("https://www.zhihu.com/creator/manage/creation/article"); err != nil {
		logger.Error("导航到页面失败", "error", err)
		return err
	}

	data, err := scraper.ExtractData(page)
	if err != nil {
		logger.Error("数据提取失败", "error", err)
		return err
	}

	logger.Info("数据提取完成", "data", data)
	return nil
}
