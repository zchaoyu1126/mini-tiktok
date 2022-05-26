package repository

import (
	"mini-tiktok/utils"
	"time"
)

type User struct {
	ID            int64     `gorm:"primary_key"`
	UserID        int64     `gorm:"Column:user_id;UNIQUE;NOT NULL"`
	UserName      string    `gorm:"Column:username;UNIQUE;NOT NULL"`
	Password      string    `gorm:"Column:password;NOT NULL"`
	CreatedAt     time.Time `gorm:"Column:created_at"`
	FollowCount   int64     `gorm:"Column:follow_count"`
	FollowerCount int64     `gorm:"Column:follower_count"`
	Status        uint8     `gorm:"Column:status;index"`
}

func (m *MysqlDao) QueryUserByName(username string) (*User, error) {
	m.db.AutoMigrate(&User{})

	user := new(User)
	err := m.db.Where("username=?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *MysqlDao) QueryUserByID(id int64) (*User, error) {
	m.db.AutoMigrate(&User{})

	user := new(User)
	err := m.db.Where("user_id=?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *MysqlDao) AddUser(username, password string) (*User, error) {
	m.db.AutoMigrate(&User{})

	user := new(User)
	snowFlake, err := utils.NewSnowFlake(0, 0)
	if err != nil {
		return nil, err
	}
	user.UserID, err = snowFlake.NextId()
	if err != nil {
		return nil, err
	}
	user.UserName = username
	user.Password = password
	user.Status = 1

	if err := m.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
