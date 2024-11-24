package config

import (
	"crawler/pkg/logger"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 总配置结构
type Config struct {
	App    AppConfig           `yaml:"app"`
	Logger logger.LoggerConfig `yaml:"logger"`
	Server Server              `yaml:"server"`
	MySQL  MySQLConfig         `yaml:"mysql"`
}

// AppConfig 应用配置结构
type AppConfig struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	CookiesFilePath string `yaml:"cookiesFilePath"`
}

// Server 服务配置
type Server struct {
	Port           string        `yaml:"port"`
	Mode           string        `yaml:"mode"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
	TrustedProxies []string      `yaml:"trustedProxies"`
	AllowedOrigins []string      `yaml:"allowedOrigins"`
	AllowedMethods []string      `yaml:"allowedMethods"`
	AllowedHeaders []string      `yaml:"allowedHeaders"`
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
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
