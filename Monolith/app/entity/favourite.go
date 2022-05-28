package entity

import "time"

type Favourite struct {
	ID          int64      `gorm:"Column:id"`
	UserID      int64      `gorm:"Column:user_id"`
	VideoID     int64      `gorm:"Column:video_id"`
	IsFavourite byte       `gorm:"Column:is_favourite"`
	CreateTime  *time.Time `gorm:"Column:create_time"`
}
