package service

import (
	"errors"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/auth"
	"mini-tiktok/common/db"
	"mini-tiktok/common/utils"

	"gorm.io/gorm"
)

type UserDTO struct {
	ID            int64
	UserName      string
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
}

// CheckUserExist query mysql to check if the username has appeared.
func CheckUserNameExist(username string) (bool, error) {
	user := &entity.User{UserName: username}
	err := dao.UserGetByName(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Before register, controller has checked wheather username exist.
// If registation success, return userID, token, nil, otherwise return -1, "", err.
func Register(username, password string) (int64, string, error) {
	// check username and password and sql injection attack
	if username == "" {
		return -1, "", errors.New("invaild username")
	} else if len(username) > 32 || len(password) > 32 {
		return -1, "", errors.New("username or password too long")
	}

	// use redis check wheather this user has already existed
	exist, err := db.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		return -1, "", err
	}
	if exist {
		return -1, "", errors.New("user already exist")
	}

	// encrypt password and then store into mysql
	encodePassword := encrypt(password)
	user := &entity.User{UserName: username, Password: encodePassword}
	snowFlake, err := utils.NewSnowFlake(0, 0)
	if err != nil {
		return -1, "", err
	}
	user.UserID, err = snowFlake.NextId()
	if err != nil {
		return -1, "", err
	}

	if err := dao.UserAdd(user); err != nil {
		return -1, "", err
	}

	// store namelist username into redis
	if err := db.NewRedisDaoInstance().AddToNameList(user.UserName); err != nil {
		return -1, "", err
	}

	// store (token, userID) into redis
	token := auth.GenerateToken(user.UserID)
	if err := db.NewRedisDaoInstance().SetToken(token, user.UserID); err != nil {
		return -1, "", err
	}
	return user.UserID, token, nil
}

// Before login, controller has checked wheather username exist.
// If login success, return userID, token, nil, otherwise return -1, "", err.
func Login(username, password string) (int64, string, error) {
	// check username and password and sql injection attack

	// use redis check wheather this user has already existed
	exist, err := db.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		return -1, "", err
	}
	if !exist {
		return -1, "", errors.New("user doesn't exist")
	}

	user := &entity.User{UserName: username}
	if err := dao.UserGetByName(user); err != nil {
		return -1, "", err
	}

	// decrypt the string stored in mysql
	decodePassword := decrypt(user.Password)
	if decodePassword != password {
		return -1, "", errors.New("wrong password")
	}

	token := auth.GenerateToken(user.UserID)
	db.NewRedisDaoInstance().SetToken(token, user.UserID)
	return user.UserID, token, nil
}

// fromUserID wants to look over toUserID's user information.
func UserInfo(toUserID, fromUserID int64) (*UserDTO, error) {
	// check toUserID and fromUserID

	user := &entity.User{UserID: toUserID}
	if err := dao.UserGetByUID(user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user id doesn't exist")
		} else {
			return nil, err
		}
	}

	// query like table, check fromUserID is follow toUserID or not
	follow := &entity.Follow{FromUserID: fromUserID, ToUserID: toUserID}
	if fromUserID == toUserID {
		follow.IsFollow = false
	} else {
		if err := dao.FollowGetByIDs(follow); err != nil {
			return nil, err
		}
	}
	userDTO := &UserDTO{user.UserID, user.UserName, user.FollowCount, user.FollowerCount, follow.IsFollow}
	return userDTO, nil
}

// use base64 encrypt password
func encrypt(str string) string {
	return str
}

// use base64 decrypt password
func decrypt(str string) string {
	return str
}
