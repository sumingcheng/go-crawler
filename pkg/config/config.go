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
}

// Server 服务配置
type Server struct {
	Port           string        `yaml:"port"`
	Mode           string        `yaml:"mode"`           // gin mode: debug/release/test
	ReadTimeout    time.Duration `yaml:"readTimeout"`    // 读取超时时间
	WriteTimeout   time.Duration `yaml:"writeTimeout"`   // 写入超时时间
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"` // 最大请求头大小
	TrustedProxies []string      `yaml:"trustedProxies"` // 受信任的代理
	AllowedOrigins []string      `yaml:"allowedOrigins"` // CORS 允许的域名
	AllowedMethods []string      `yaml:"allowedMethods"` // CORS 允许的方法
	AllowedHeaders []string      `yaml:"allowedHeaders"` // CORS 允许的请求头
}

// AppConfig 应用配置结构
type AppConfig struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	CookiesFilePath string `yaml:"cookiesFilePath"`
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
