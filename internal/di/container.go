package di

import (
	"crawler/internal/controller"
	"crawler/internal/repository"
	"crawler/internal/router"
	"crawler/internal/service"
	"crawler/pkg/config"

	"gorm.io/gorm"
)

type Container struct {
	Config         *config.Config
	DB             *gorm.DB
	ArticleRepo    repository.ArticleRepository
	CrawlerService service.ICrawlerService
	CrawlerHandler *controller.CrawlerController
	Router         router.IRouter
}

func NewContainer(cfg *config.Config, db *gorm.DB) (*Container, error) {
	container := &Container{
		Config: cfg,
		DB:     db,
	}

	if err := container.initializeDependencies(); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initializeDependencies() error {
	// 1. Repository
	c.ArticleRepo = repository.NewGormArticleRepository(c.DB)

	// 2. Service
	c.CrawlerService = service.NewCrawlerService(c.Config, c.ArticleRepo)

	// 3. Controller
	c.CrawlerHandler = controller.NewCrawlerController(c.CrawlerService)

	// 4. Router
	r, err := router.NewRouter(c.Config, controller.NewHandlers(c.CrawlerHandler))
	if err != nil {
		return err
	}
	c.Router = r

	return nil
}

// 添加清理方法
func (c *Container) Cleanup() {
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
