package config

import (
	"crawler/pkg/logger"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 总配置结构
type Config struct {
	App    AppConfig           `yaml:"app"`
	Logger logger.LoggerConfig `yaml:"logger"`
}

// AppConfig 应用配置结构
type AppConfig struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	CookiesFilePath string `yaml:"cookiesFilePath"`
	// 可以根据需要添加其他配置项
}

// LoadConfig 加载配置文件
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
