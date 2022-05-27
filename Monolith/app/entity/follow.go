package entity

import "time"

type Follow struct {
	ID         int64 `gorm:"Column:id"`
	FromUserID int64 `gorm:"Column:from_user_id"`
	ToUserID   int64 `gorm:"Column:to_user_id"`
	IsFollow   bool  `gorm:"is_follow"`
	CreatedAt  *time.Time
	DeletedAt  *time.Time
}
