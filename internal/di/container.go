package di

import (
	"crawler/internal/controller"
	"crawler/internal/repository"
	"crawler/internal/router"
	"crawler/internal/service"
	"crawler/pkg/config"
	"fmt"

	"gorm.io/gorm"
)

type Container struct {
	Config         *config.Config
	DB             *gorm.DB
	ArticleRepo    repository.ArticleRepository
	CrawlerService service.ICrawlerService
	CrawlerHandler controller.ICrawlerController
	Router         *router.Router
}

func NewContainer(cfg *config.Config, db *gorm.DB) (*Container, error) {
	// 1. Repository
	articleRepo := repository.NewGormArticleRepository(db)

	// 2. Service
	crawlerService := service.NewCrawlerService(cfg, articleRepo)

	// 3. Controller
	crawlerController := controller.NewCrawlerController(crawlerService)

	// 4. Router
	r, err := router.NewRouter(cfg, crawlerController)
	if err != nil {
		return nil, fmt.Errorf("初始化路由失败: %w", err)
	}

	return &Container{
		Config:         cfg,
		DB:             db,
		ArticleRepo:    articleRepo,
		CrawlerService: crawlerService,
		CrawlerHandler: crawlerController,
		Router:         r,
	}, nil
}

// 添加清理方法
func (c *Container) ReleaseResources() {
	// 按依赖关系的反向顺序清理资源
	if c.CrawlerService != nil {
		c.CrawlerService.Cleanup()
	}

	if c.DB != nil {
		if sqlDB, err := c.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}
