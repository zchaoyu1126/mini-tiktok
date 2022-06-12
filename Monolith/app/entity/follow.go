package entity

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	ID         int64 `gorm:"Column:id;primary_key"`
	FromUserID int64 `gorm:"Column:from_user_id"`
	ToUserID   int64 `gorm:"Column:to_user_id"`
	IsFollow   bool  `gorm:"is_follow"`
}
