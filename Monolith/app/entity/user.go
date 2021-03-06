package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            int64  `gorm:"Column:id;primary_key"`
	UserID        int64  `gorm:"Column:user_id;UNIQUE;NOT NULL"`
	UserName      string `gorm:"Column:username;UNIQUE;NOT NULL"`
	Password      string `gorm:"Column:password;NOT NULL"`
	FollowCount   int64  `gorm:"Column:follow_count"`
	FollowerCount int64  `gorm:"Column:follower_count"`
}
