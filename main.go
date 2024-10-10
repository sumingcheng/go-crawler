package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type Config struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	CookiesFilePath string `yaml:"cookiesFilePath"`
}

func main() {
	err, username, password, cookiesFilePath := ConfigInit()

	// 初始化 Playwright
	if err := playwright.Install(); err != nil {
		log.Fatalf("无法安装 Playwright: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("无法启动 Playwright: %v", err)
	}
	defer pw.Stop()

	// 启动浏览器
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // 设置为 true 则为无头模式
	})
	if err != nil {
		log.Fatalf("无法启动浏览器: %v", err)
	}
	defer browser.Close()

	// 初始化上下文
	var context playwright.BrowserContext

	if _, err := os.ReadFile(cookiesFilePath); err == nil {
		// 如果存在 storage_state.json，加载存储状态
		context, err = browser.NewContext(playwright.BrowserNewContextOptions{
			StorageStatePath: playwright.String(cookiesFilePath),
		})
		if err != nil {
			log.Fatalf("无法创建浏览器上下文: %v", err)
		}
	} else {
		// 如果不存在 storage_state.json，创建新上下文并登录
		context, err = browser.NewContext()
		if err != nil {
			log.Fatalf("无法创建浏览器上下文: %v", err)
		}

		page, err := context.NewPage()
		if err != nil {
			log.Fatalf("无法创建页面: %v", err)
		}

		// 导航到登录页面
		if _, err = page.Goto("https://www.zhihu.com/signin"); err != nil {
			log.Fatalf("无法导航到登录页面: %v", err)
		}

		// 输入用户名和密码（请替换为您的实际用户名和密码）
		if err = page.Fill("input[name='username']", username); err != nil {
			log.Fatalf("无法输入用户名: %v", err)
		}

		if err = page.Fill("input[name='password']", password); err != nil {
			log.Fatalf("无法输入密码: %v", err)
		}

		// 点击登录按钮
		if err = page.Click("button[type='submit']"); err != nil {
			log.Fatalf("无法点击登录按钮: %v", err)
		}

		// 等待登录成功（替换为登录后页面的特定元素选择器）
		if _, err = page.WaitForSelector("img[alt='头像']"); err != nil {
			log.Fatalf("登录后未能成功导航: %v", err)
		}

		// 获取存储状态
		state, err := context.StorageState()
		if err != nil {
			log.Fatalf("无法获取存储状态: %v", err)
		}

		// 序列化存储状态为 JSON
		stateJSON, err := json.Marshal(state)
		if err != nil {
			log.Fatalf("无法序列化存储状态: %v", err)
		}

		// 保存存储状态到文件
		if err = ioutil.WriteFile(cookiesFilePath, stateJSON, 0644); err != nil {
			log.Fatalf("无法保存存储状态: %v", err)
		}
	}
	defer context.Close()

	// 创建新的页面
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("无法创建页面: %v", err)
	}

	// 监听网络请求响应
	var apiResponse playwright.Response
	page.OnResponse(func(response playwright.Response) {
		if strings.Contains(response.URL(), "/api/v4/creators/creations/v2/article") {
			apiResponse = response
		}
	})

	// 导航到需要触发 API 请求的页面
	if _, err := page.Goto("https://www.zhihu.com/creator/manage/creation/article"); err != nil {
		log.Fatalf("无法导航到页面: %v", err)
	}

	// 延迟一段时间以确保请求被捕获到
	// 这里可以设置更长时间，也可以使用 page.WaitForTimeout(xxx) 来确保加载完成
	if apiResponse != nil {
		data, err := apiResponse.Text()
		if err != nil {
			log.Fatalf("无法获取响应内容: %v", err)
		}

		// 检查是否获取到数据
		if data != "" {
			fmt.Println("获取到的 API 数据：")
			fmt.Println(data)
		} else {
			fmt.Println("未能获取到 API 数据，请检查是否已登录或请求 URL 是否正确。")
		}
	} else {
		fmt.Println("未能获取到目标 API 请求的响应。")
	}

}

func ConfigInit() (error, string, string, string) {
	// 读取 YAML 配置文件
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	// 解析 YAML 配置文件
	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}

	// 获取用户名和密码
	username := config.Username
	password := config.Password
	// 获取存储状态文件路径
	cookiesFilePath := config.CookiesFilePath

	if username == "" || password == "" {
		log.Fatalf("用户名或密码为空，请检查配置文件")
	}
	return err, username, password, cookiesFilePath
}
