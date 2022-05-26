package service

import (
	"errors"
	"fmt"
	"mini-tiktok/repository"
	"mini-tiktok/utils"

	"gorm.io/gorm"
)

type UserView struct {
	ID            int64  `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// CheckUserExist query mysql to check if the username has appeared.
func CheckUserExist(username string) (bool, error) {
	user, err := repository.NewMysqlDaoInstance().QueryUserByName(username)
	if err != nil {
		return false, err
	} else if user == nil {
		return false, nil
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

	// encrypt password and then store into mysql
	encodePassword := encrypt(password)
	user, err := repository.NewMysqlDaoInstance().AddUser(username, encodePassword)
	if err != nil {
		return -1, "", err
	}

	// store namelist username into redis
	err = repository.NewRedisDaoInstance().AddToNameList(user.UserName)
	if err != nil {
		return -1, "", err
	}

	// store (token, userID) into redis
	token := utils.GenerateToken(user.UserID)
	err = repository.NewRedisDaoInstance().SetToken(token, user.UserID)
	if err != nil {
		return -1, "", err
	}
	return user.UserID, token, nil
}

// Before login, controller has checked wheather username exist.
// If login success, return userID, token, nil, otherwise return -1, "", err.
func Login(username, password string) (int64, string, error) {
	// check username and password and sql injection attack

	user, err := repository.NewMysqlDaoInstance().QueryUserByName(username)
	if err != nil {
		return -1, "", err
	}

	// decrypt the string stored in mysql
	decodePassword := decrypt(user.Password)
	if decodePassword != password {
		return -1, "", errors.New("wrong password")
	}

	token := utils.GenerateToken(user.UserID)
	repository.NewRedisDaoInstance().SetToken(token, user.UserID)
	return user.UserID, token, nil
}

// fromUserID wants to look over toUserID's user information.
func UserInfo(toUserID, fromUserID int64) (*UserView, error) {
	// check toUserID and fromUserID

	userView := new(UserView)
	user, err := repository.NewMysqlDaoInstance().QueryUserByID(toUserID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user id doesn't exist")
	}
	if err != nil {
		return nil, err
	}

	// query like table, check fromUserID is follow toUserID or not
	var isFollow bool
	fmt.Println(fromUserID, toUserID, fromUserID == toUserID)
	if fromUserID == toUserID {
		isFollow = false
	} else {
		isFollow, err = repository.NewMysqlDaoInstance().IsFollow(fromUserID, toUserID)
	}

	if err != nil {
		return nil, err
	}

	userView.ID = user.UserID
	userView.UserName = user.UserName
	userView.FollowCount = user.FollowCount
	userView.FollowerCount = user.FollowerCount
	userView.IsFollow = isFollow
	return userView, nil
}

// use base64 encrypt password
func encrypt(str string) string {
	return str
}

// use base64 decrypt password
func decrypt(str string) string {
	return str
}
