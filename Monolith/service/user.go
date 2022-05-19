package service

import (
	"errors"
	"fmt"
	"mini-tiktok/middleware"
	"mini-tiktok/repository"
)

type UserInformation struct {
	ID            int64
	UserName      string
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
}

func CheckUserExist(username string) (bool, error) {
	user, err := repository.NewDaoInstance().QueryUserByName(username)
	if err != nil {
		return false, err
	} else if user == nil {
		return false, nil
	}
	return true, nil
}

func Register(username, password string) (int64, string, error) {
	// 需要对password使用哈希加密，不能直接存password
	// 之后完善
	if username == "" {
		return -1, "", errors.New("invaild username")
	}
	encodePassword := Encryption(password)
	user, err := repository.NewDaoInstance().AddUser(username, encodePassword)
	if err != nil {
		return -1, "", err
	}
	claims := middleware.MyCustomClaims{ID: user.ID, Username: user.UserName}
	token, _ := middleware.GenerateToken(claims)
	return user.ID, token, nil
}

// Login之前已经判断过是否存在该用户
func Login(username, password string) (int64, string, error) {
	user, err := repository.NewDaoInstance().QueryUserByName(username)
	if err != nil {
		return -1, "", err
	}
	decodePassword := Decryption(user.Password)
	if decodePassword != password {
		return -1, "", errors.New("wrong password")
	}
	claims := middleware.MyCustomClaims{ID: user.ID, Username: user.UserName}
	token, _ := middleware.GenerateToken(claims)
	return user.ID, token, nil
}

func UserInfo(id, callerID int64) (*UserInformation, error) {
	// 检查id是否合法
	userInformation := new(UserInformation)
	if id < 0 {
		return nil, errors.New("invaild id")
	}

	// 向数据库中查询ID和用户名
	fmt.Println(id, callerID)
	user, err := repository.NewDaoInstance().QueryUserByID(id)

	if user == nil && err == nil {
		return nil, errors.New("user id doesn't exist")
	}
	fmt.Println(id, callerID, user)
	userInformation.ID = user.ID
	userInformation.UserName = user.UserName
	// 向数据库中查询follow_count
	userInformation.FollowCount = 0
	// 向数据库中查询follower_count
	userInformation.FollowerCount = 0
	// 向数据库中查询id与callerID之间的关系
	userInformation.IsFollow = false
	return userInformation, nil
}

func Encryption(str string) string {
	return str
}

func Decryption(str string) string {
	return str
}
