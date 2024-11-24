package repository

import (
	"crawler/internal/repository/model"
	"crawler/internal/scraper"
	"crawler/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleRepository interface {
	UpsertArticles(articles []scraper.ArticleCard) error
	FindAll() ([]scraper.ArticleCard, error)
}

type GormArticleRepository struct {
	db *gorm.DB
}

func NewGormArticleRepository(db *gorm.DB) ArticleRepository {
	return &GormArticleRepository{db: db}
}

func (r *GormArticleRepository) UpsertArticles(articles []scraper.ArticleCard) error {
	// 将爬虫数据转换为数据库模型
	var models []model.Article
	for _, article := range articles {
		models = append(models, model.Article{
			Title:         article.Title,
			Link:          article.Link,
			Description:   article.Description,
			PublishedTime: article.PublishedTime,
			ViewCount:     article.Stats.Reads,
			Upvote:        article.Stats.Upvote,
			Comments:      article.Stats.Comments,
			Bookmarks:     article.Stats.Bookmarks,
			Likes:         article.Stats.Likes,
			Status:        1,
		})
	}

	// 使用 Upsert 进行批量插入或更新
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "link"}},
		DoUpdates: clause.AssignmentColumns([]string{"view_count", "upvote", "comments", "bookmarks", "likes"}),
	}).Create(&models)

	if result.Error != nil {
		logger.Error("保存文章失败", "error", result.Error)
		return result.Error
	}

	logger.Info("成功保存文章", "count", len(articles))
	return nil
}

func (r *GormArticleRepository) FindAll() ([]scraper.ArticleCard, error) {
	var articles []model.Article
	if err := r.db.Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, err
	}

	// 转换为爬虫数据结构
	result := make([]scraper.ArticleCard, len(articles))
	for i, article := range articles {
		result[i] = scraper.ArticleCard{
			Title:         article.Title,
			Link:          article.Link,
			Description:   article.Description,
			PublishedTime: article.PublishedTime,
			Stats: scraper.ArticleStats{
				Reads:     article.ViewCount,
				Upvote:    article.Upvote,
				Comments:  article.Comments,
				Bookmarks: article.Bookmarks,
				Likes:     article.Likes,
			},
		}
	}

	return result, nil
}
