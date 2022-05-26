package repository

import "time"

type Follow struct {
	ID         int64     `gorm:"Column:id"`
	FromUserID int64     `gorm:"Column:from_user_id"`
	ToUserID   int64     `gorm:"Column:to_user_id"`
	IsFollow   bool      `gorm:"is_follow"`
	CreatedAt  time.Time `gorm:"Column:created_at"`
}

func (m *MysqlDao) IsFollow(fromUserID, toUserID int64) (bool, error) {
	m.db.AutoMigrate(&Follow{})

	follow := new(Follow)
	err := m.db.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).First(follow).Error
	if err != nil {
		return false, err
	}
	return follow.IsFollow, nil
}

func (m *MysqlDao) Follow(fromUserID, toUserID int64) error {
	m.db.AutoMigrate(&Follow{})

	follow := new(Follow)
	follow.FromUserID = fromUserID
	follow.ToUserID = toUserID
	follow.IsFollow = true
	// 首先查询是否存在
	// 如果不存在，那么直接create
	// 如果存在，则修改
	return m.db.Create(follow).Error
}

func (m *MysqlDao) DisFollow(fromUserID, toUserID int64) error {
	m.db.AutoMigrate(&Follow{})

	follow := new(Follow)
	m.db.Model(follow).Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Update("is_follow", false)
	return nil
}

func (m *MysqlDao) ModifyFollowCount(userID, newVal int64) error {
	m.db.AutoMigrate(&User{})

	user := new(User)
	m.db.Model(user).Where("user_id = ?", userID).Update("follow_count", newVal)
	return nil
}

func (m *MysqlDao) ModifyFollowerCount(userID, newVal int64) error {
	m.db.AutoMigrate(&User{})

	user := new(User)
	m.db.Model(user).Where("user_id = ?", userID).Update("follower_count", newVal)
	return nil
}

func (m *MysqlDao) UpdateUser(user *User) error {
	m.db.Save(user)
	return nil
}
