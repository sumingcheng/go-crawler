package cookies

import (
	"crawler/pkg/logger"
	"encoding/json"
	"fmt"
	"os"

	"github.com/playwright-community/playwright-go"
)

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

func LoadCookies(context playwright.BrowserContext, cookiesFilePath string) error {
	logger.Info("开始加载Cookies",
		"file_path", cookiesFilePath,
	)

	cookiesData, err := os.ReadFile(cookiesFilePath)
	if err != nil {
		logger.Error("读取Cookies文件失败",
			"error", err,
			"file_path", cookiesFilePath,
		)
		return fmt.Errorf("读取Cookies文件失败: %w", err)
	}

	var originalCookies []OriginalCookie
	if err := json.Unmarshal(cookiesData, &originalCookies); err != nil {
		logger.Error("解析Cookies数据失败",
			"error", err,
		)
		return fmt.Errorf("解析Cookies数据失败: %w", err)
	}

	logger.Info("成功解析Cookies数据",
		"count", len(originalCookies),
	)

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
		logger.Error("添加Cookies到浏览器失败",
			"error", err,
			"cookies_count", len(cookies),
		)
		return fmt.Errorf("添加Cookies失败: %w", err)
	}

	logger.Info("Cookies添加成功",
		"count", len(cookies),
	)
	return nil
}
