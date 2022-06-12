package entity

import (
	"time"

	"gorm.io/gorm"
)

type Publication struct {
	gorm.Model
	ID             int64      `gorm:"Column:id;primary_key"`
	VideoID        int64      `gorm:"Column:video_id;UNIQUE"`
	OwnerID        int64      `gorm:"Column:owner_id"`
	Title          string     `gorm:"Column:title"`
	PlayUrl        string     `gorm:"Column:play_url"`
	CoverUrl       string     `gorm:"Column:cover_url"`
	CreateTime     *time.Time `gorm:"Column:create_time"`
	FavouriteCount int        `gorm:"Column:favourite_count"`
	CommentCount   int        `gorm:"Column:comment_count"`
}
