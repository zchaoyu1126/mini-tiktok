package entity

import (
	"time"
)

type Publication struct {
	ID             int64      `gorm:"primary_key"`
	VideoID        int64      `gorm:"Column:video_id;UNIQUE"`
	OwnerID        int64      `gorm:"Column:owner_id"`
	Title          string     `gorm:"Column:title"`
	PlayUrl        string     `gorm:"Column:play_url"`
	CoverUrl       string     `gorm:"Column:cover_url"`
	CreateTime     *time.Time `gorm:"Column:create_time"`
	FavouriteCount int        `gorm:"Column:favourite_count"`
	CommentCount   int        `gorm:"Column:comment_count"`
	Status         int16      `gorm:"Column:status"`
}
