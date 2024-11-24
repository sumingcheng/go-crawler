package repository

const (
	CreateDatabaseSQL = `
			CREATE DATABASE IF NOT EXISTS %s
			DEFAULT CHARACTER SET utf8mb4
			DEFAULT COLLATE utf8mb4_unicode_ci;
	`

	CreateArticlesTableSQL = `
			CREATE TABLE IF NOT EXISTS articles (
					id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
					title VARCHAR(255) NOT NULL COMMENT '文章标题',
					link VARCHAR(512) NOT NULL COMMENT '文章链接',
					description TEXT COMMENT '文章描述',
					published_time VARCHAR(64) COMMENT '发布时间',
					view_count INT UNSIGNED DEFAULT 0 COMMENT '阅读数',
					upvote INT UNSIGNED DEFAULT 0 COMMENT '点赞数',
					comments INT UNSIGNED DEFAULT 0 COMMENT '评论数',
					bookmarks INT UNSIGNED DEFAULT 0 COMMENT '收藏数',
					likes INT UNSIGNED DEFAULT 0 COMMENT '喜欢数',
					status TINYINT(1) DEFAULT 1 COMMENT '状态:1-正常,0-删除',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
					PRIMARY KEY (id),
					UNIQUE KEY uk_link (link)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';
	`
)
