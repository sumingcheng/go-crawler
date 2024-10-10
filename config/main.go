package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CookiesFilePath string `yaml:"cookies_file_path"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	configFile, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return config, fmt.Errorf("解析配置文件失败: %v", err)
	}

	if config.CookiesFilePath == "" {
		return config, fmt.Errorf("cookies 文件路径为空，请检查配置文件")
	}

	return config, nil
}
