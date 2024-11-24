package repository

import (
	"crawler/internal/scraper"
	"crawler/pkg/logger"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticleRepository interface {
	Save(articles []scraper.ArticleCard) error
	FindAll() ([]scraper.ArticleCard, error)
}

type MySQLArticleRepository struct {
	db *sqlx.DB
}

func NewMySQLArticleRepository(db *sqlx.DB) ArticleRepository {
	return &MySQLArticleRepository{db: db}
}

type ArticleModel struct {
	ID            int64     `db:"id"`
	Title         string    `db:"title"`
	Link          string    `db:"link"`
	Description   string    `db:"description"`
	PublishedTime string    `db:"published_time"`
	Reads         int       `db:"reads"`
	Upvote        int       `db:"upvote"`
	Comments      int       `db:"comments"`
	Bookmarks     int       `db:"bookmarks"`
	Likes         int       `db:"likes"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (r *MySQLArticleRepository) Save(articles []scraper.ArticleCard) error {
	query := `
        INSERT INTO articles (
            title, link, description, published_time, 
            reads, upvote, comments, bookmarks, likes, created_at
        ) VALUES (
            :title, :link, :description, :published_time,
            :reads, :upvote, :comments, :bookmarks, :likes, :created_at
        ) ON DUPLICATE KEY UPDATE 
            reads=VALUES(reads),
            upvote=VALUES(upvote),
            comments=VALUES(comments),
            bookmarks=VALUES(bookmarks),
            likes=VALUES(likes),
            updated_at=CURRENT_TIMESTAMP
    `

	for _, article := range articles {
		model := ArticleModel{
			Title:         article.Title,
			Link:          article.Link,
			Description:   article.Description,
			PublishedTime: article.PublishedTime,
			Reads:         article.Stats.Reads,
			Upvote:        article.Stats.Upvote,
			Comments:      article.Stats.Comments,
			Bookmarks:     article.Stats.Bookmarks,
			Likes:         article.Stats.Likes,
			CreatedAt:     time.Now(),
		}

		_, err := r.db.NamedExec(query, model)
		if err != nil {
			logger.Error("保存文章失败",
				"error", err,
				"title", article.Title,
				"link", article.Link,
			)
			return fmt.Errorf("failed to save article: %w", err)
		}
	}

	logger.Info("成功保存文章", "count", len(articles))
	return nil
}

func (r *MySQLArticleRepository) FindAll() ([]scraper.ArticleCard, error) {
	query := `
        SELECT title, link, description, published_time,
               reads, upvote, comments, bookmarks, likes
        FROM articles ORDER BY created_at DESC
    `

	var models []ArticleModel
	if err := r.db.Select(&models, query); err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}

	articles := make([]scraper.ArticleCard, len(models))
	for i, model := range models {
		articles[i] = scraper.ArticleCard{
			Title:         model.Title,
			Link:          model.Link,
			Description:   model.Description,
			PublishedTime: model.PublishedTime,
			Stats: scraper.ArticleStats{
				Reads:     model.Reads,
				Upvote:    model.Upvote,
				Comments:  model.Comments,
				Bookmarks: model.Bookmarks,
				Likes:     model.Likes,
			},
		}
	}

	return articles, nil
}
