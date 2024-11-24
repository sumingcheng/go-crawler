package repository

import "time"

// Article GORM 文章模型
type Article struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Title         string    `gorm:"type:varchar(255);not null;comment:文章标题"`
	Link          string    `gorm:"type:varchar(512);not null;uniqueIndex:uk_link;comment:文章链接"`
	Description   string    `gorm:"type:text;comment:文章描述"`
	PublishedTime string    `gorm:"type:varchar(64);comment:发布时间"`
	ViewCount     int       `gorm:"type:int unsigned;default:0;comment:阅读数"`
	Upvote        int       `gorm:"type:int unsigned;default:0;comment:点赞数"`
	Comments      int       `gorm:"type:int unsigned;default:0;comment:评论数"`
	Bookmarks     int       `gorm:"type:int unsigned;default:0;comment:收藏数"`
	Likes         int       `gorm:"type:int unsigned;default:0;comment:喜欢数"`
	Status        int8      `gorm:"type:tinyint(1);default:1;comment:状态:1-正常,0-删除"`
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}
