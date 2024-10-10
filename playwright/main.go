package playwright

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func InitializePlaywright() (*playwright.Playwright, playwright.Browser, error) {
	log.Println("开始安装 Playwright")
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Playwright 启动成功")

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		pw.Stop()
		return nil, nil, err
	}
	log.Println("浏览器启动成功")

	return pw, browser, nil
}
