package service

import (
	"errors"
	"mini-tiktok/app/dao"
	"mini-tiktok/app/entity"
	"mini-tiktok/common/auth"
	"mini-tiktok/common/db"
	"mini-tiktok/common/utils"
	"mini-tiktok/common/xerr"
	"regexp"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserVO struct {
	ID            int64  `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// 去数据库中查询用户名是否存在
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

func Register(username, password string) (int64, string, error) {
	// 使用正则表达式检验用户输入的用户名和密码是否合法，
	// 不合法则返回ErrUsernameValidation或ErrPasswordValidation

	// 用户名由字母数据下划线英文句号组成，长度要求4-16之间
	usernameReg, _ := regexp.Compile(`^[a-zA-Z0-9_\.@]{4,16}$`)
	if !usernameReg.MatchString(username) {
		return 0, "", xerr.ErrUsernameValidation
	}

	// 密码匹配6-16位英文数据大部分英文标点
	passwordReg, _ := regexp.Compile(`^([A-Za-z0-9\-=\[\];,\./~!@#\$%^\*\(\)_\+}{:\?]){6,16}$`)
	if !passwordReg.MatchString(password) {
		return 0, "", xerr.ErrPasswordValidation
	}

	// 使用redis查询用户名是否已存在，如果已存在返回ErrUserExist
	exist, err := db.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		return 0, "", err
	}
	if exist {
		return 0, "", xerr.ErrUserExist
	}

	// 使用标准库函数对用户的密码进行加密
	encodePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("encrypt password failed")
		return 0, "", xerr.ErrInternalServer
	}
	user := &entity.User{UserName: username, Password: string(encodePassword)}

	// 初始化雪花算法，并生成userid
	snowFlake, err := utils.NewSnowFlake(0, 0)
	if err != nil {
		zap.L().Error("init snow flake failed")
		return 0, "", xerr.ErrInternalServer
	}
	user.UserID, err = snowFlake.NextId()
	if err != nil {
		zap.L().Error("generate uid failed")
		return 0, "", xerr.ErrInternalServer
	}

	// 将用户信息存储在mysql中
	if err := dao.UserAdd(user); err != nil {
		return 0, "", err
	}

	// 将用户名信息存储在redis中
	if err := db.NewRedisDaoInstance().AddToNameList(user.UserName); err != nil {
		return 0, "", err
	}

	// 将token,userid 存储在redis中
	token := auth.GenerateToken(user.UserID)
	if err := db.NewRedisDaoInstance().SetToken(token, user.UserID); err != nil {
		return 0, "", err
	}
	return user.UserID, token, nil
}

func Login(username, password string) (int64, string, error) {
	// 使用redis查询用户名是否存在，如果不存在则直接返回ErrUserNotFound错误
	exist, err := db.NewRedisDaoInstance().IsUserNameExist(username)
	if err != nil {
		return 0, "", err
	}
	if !exist {
		return 0, "", xerr.ErrUserNotFound
	}

	// 查询mysql数据库，因为已经在redis中查询过，所以不需要考虑用户是否存在的问题
	// 即使不存在，err为gorm.ErrRecordNotFound那么也会被直接返回
	user := &entity.User{UserName: username}
	if err := dao.UserGetByName(user); err != nil {
		return 0, "", err
	}

	// 将用户输入的password与数据库中存储的密文比较，比如密码不正确，则返回ErrPasswordIncorrect
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, "", xerr.ErrPasswordIncorrect
	}

	// 将token, user_id键值对存储在redis中
	token := auth.GenerateToken(user.UserID)
	if err = db.NewRedisDaoInstance().SetToken(token, user.UserID); err != nil {
		return 0, "", err
	}
	return user.UserID, token, nil
}

// fromUserID wants to look over toUserID's user information.
func UserInfo(toUserID, fromUserID int64) (*UserVO, error) {

	user := &entity.User{UserID: toUserID}
	// 查询想查看的用户是否存在
	if err := dao.UserGetByUID(user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerr.ErrUserNotFound
		}
		return nil, err
	}

	// query like table, check fromUserID is follow toUserID or not
	follow := &entity.Follow{FromUserID: fromUserID, ToUserID: toUserID}
	var isFollow bool
	if fromUserID == toUserID {
		isFollow = false
	} else {
		if err := dao.FollowGetByIDs(follow); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isFollow = false
			} else if err != nil {
				return nil, err
			}
		}
		isFollow = follow.IsFollow
	}
	userDTO := &UserVO{user.UserID, user.UserName, user.FollowCount, user.FollowerCount, isFollow}
	return userDTO, nil
}
