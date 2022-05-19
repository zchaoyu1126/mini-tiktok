package repository

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int64  `gorm:"primary_key"`
	UserName string `gorm:"Column:username;UNIQUE;NOT NULL"`
	Password string `gorm:"Column:password;NOT NULL"`
}

func (m *MysqlDao) QueryUserByName(username string) (*User, error) {
	m.db.AutoMigrate(&User{})
	// new(T) 分配了零值填充的 T 类型的内存空间，并且返回其地址，一个 *T 类型的值
	user := new(User)

	err := m.db.Where("username=?", username).First(user).Error
	if err == nil {
		// 在数据库中找到该用户
		return user, nil
	} else {
		// 错误为未查询到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			// 数据库过程查询出错
			log.Println(err)
			return nil, err
		}
	}
}

func (m *MysqlDao) QueryUserByID(id int64) (*User, error) {
	m.db.AutoMigrate(&User{})
	// new(T) 分配了零值填充的 T 类型的内存空间，并且返回其地址，一个 *T 类型的值
	user := new(User)

	err := m.db.Where("id=?", id).First(user).Error
	if err == nil {
		// 在数据库中找到该用户
		return user, nil
	} else {
		// 错误为未查询到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			// 数据库过程查询出错
			log.Println(err)
			return nil, err
		}
	}
}

func (m *MysqlDao) AddUser(username, password string) (*User, error) {
	m.db.AutoMigrate(&User{})
	user := new(User)
	user.UserName = username
	user.Password = password
	if err := m.db.Create(user).Error; err != nil {
		// 数据库写数据时出错
		log.Println(err)
		return nil, err
	}
	return user, nil
}
