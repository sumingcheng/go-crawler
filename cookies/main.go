package cookies

import (
	"encoding/json"
	"fmt"
	"log"
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
