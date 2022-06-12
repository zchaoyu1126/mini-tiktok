package entity

import (
	"gorm.io/gorm"
)

type Favourite struct {
	gorm.Model
	ID          int64 `gorm:"Column:id;primary_key"`
	UserID      int64 `gorm:"Column:user_id"`
	VideoID     int64 `gorm:"Column:video_id"`
	IsFavourite bool  `gorm:"Column:is_favourite"`
}
