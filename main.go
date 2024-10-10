package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/playwright-community/playwright-go"
	"gopkg.in/yaml.v3"
)

// Config 定义配置文件结构
type Config struct {
	CookiesFilePath string `yaml:"cookies_file_path"`
}

// OriginalCookie 定义原始Cookie结构
type OriginalCookie struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate,omitempty"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	HttpOnly       bool    `json:"httpOnly"`
	Secure         bool    `json:"secure"`
	Value          string  `json:"value"`
}

func main() {
	log.Println("开始执行程序")

	// 安装Playwright
	if err := playwright.Install(); err != nil {
		log.Fatalf("安装 Playwright 失败: %v", err)
	} else {
		log.Println("Playwright 安装成功")
	}

	// 运行Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("启动 Playwright 失败: %v", err)
	}
	defer func() {
		pw.Stop()
		log.Println("Playwright 停止运行")
	}()

	// 启动浏览器
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("启动浏览器失败: %v", err)
	}
	defer func() {
		browser.Close()
		log.Println("浏览器已关闭")
	}()

	// 初始化配置
	cookiesFilePath, err := ConfigInit()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 检查Cookies文件路径
	if _, err := os.Stat(cookiesFilePath); os.IsNotExist(err) {
		log.Fatalf("Cookies 文件不存在，请检查路径: %v", cookiesFilePath)
	} else {
		log.Printf("Cookies 文件路径: %s", cookiesFilePath)
	}

	// 创建浏览器上下文
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("创建浏览器上下文失败: %v", err)
	}
	defer func() {
		context.Close()
		log.Println("浏览器上下文已关闭")
	}()

	// 加载Cookies
	err = LoadCookies(context, cookiesFilePath)
	if err != nil {
		log.Fatalf("加载Cookies失败: %v", err)
	} else {
		log.Println("Cookies 加载成功")
	}

	// 创建页面
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("创建页面失败: %v", err)
	}
	defer func() {
		page.Close()
		log.Println("页面已关闭")
	}()

	// 导航到页面
	if _, err := page.Goto("https://www.zhihu.com/creator/manage/creation/article"); err != nil {
		log.Fatalf("导航到页面失败: %v", err)
	} else {
		log.Println("成功导航到指定页面")
	}

	// 等待响应
	timeoutOption := playwright.PageExpectResponseOptions{Timeout: playwright.Float(60000)} // 增加超时时间
	_, err = page.ExpectResponse("api/v4/creators/creations/v2/article", func() error {
		return nil
	}, timeoutOption)
	if err != nil {
		log.Fatalf("等待响应失败: %v", err)
	} else {
		log.Println("成功获取响应")
	}
	fmt.Println("成功获取响应")
}

// ConfigInit 从配置文件初始化配置
func ConfigInit() (string, error) {
	log.Println("开始读取配置文件")
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return "", fmt.Errorf("解析配置文件失败: %v", err)
	}

	if config.CookiesFilePath == "" {
		return "", fmt.Errorf("cookies 文件路径为空，请检查配置文件")
	}
	log.Println("配置文件读取并解析成功")
	return config.CookiesFilePath, nil
}

// LoadCookies 从文件加载Cookies到浏览器上下文
func LoadCookies(context playwright.BrowserContext, cookiesFilePath string) error {
	log.Printf("开始从文件加载Cookies，文件路径: %s", cookiesFilePath)
	cookiesData, err := os.ReadFile(cookiesFilePath)
	if err != nil {
		return fmt.Errorf("读取Cookies文件失败: %v", err)
	}

	var originalCookies []OriginalCookie
	if err := json.Unmarshal(cookiesData, &originalCookies); err != nil {
		return fmt.Errorf("解析Cookies数据失败: %v", err)
	}

	cookies := make([]playwright.OptionalCookie, 0, len(originalCookies))
	for _, oc := range originalCookies {
		var samesite *playwright.SameSiteAttribute
		switch oc.SameSite {
		case "Strict":
			samesite = playwright.SameSiteAttributeStrict
		case "Lax":
			samesite = playwright.SameSiteAttributeLax
		case "None":
			samesite = playwright.SameSiteAttributeNone
		}

		var exp *float64
		if oc.ExpirationDate != 0 {
			exp = &oc.ExpirationDate
		}

		cookie := playwright.OptionalCookie{
			Name:     oc.Name,
			Value:    oc.Value,
			Domain:   &oc.Domain,
			Path:     &oc.Path,
			HttpOnly: &oc.HttpOnly,
			Secure:   &oc.Secure,
			SameSite: samesite,
			Expires:  exp,
		}
		cookies = append(cookies, cookie)
	}

	if err := context.AddCookies(cookies); err != nil {
		return fmt.Errorf("添加Cookies失败: %v", err)
	}
	log.Println("Cookies成功添加到浏览器上下文")
	return nil
}
