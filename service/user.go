package service

import (
	"errors"
	"mini-tiktok/middleware"
	"mini-tiktok/repository"
)

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
	encodePassword := Encryption(password)
	user, err := repository.NewDaoInstance().AddUser(username, encodePassword)
	if err != nil {
		return -1, "", err
	}
	claims := middleware.MyCustomClaims{Id: user.ID, Username: user.UserName}
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
	claims := middleware.MyCustomClaims{Id: user.ID, Username: user.UserName}
	token, _ := middleware.GenerateToken(claims)
	return user.ID, token, nil
}

// func GetUserInfo(id int64) *model.User {
// 	user := new(model.User)
// 	err := dao.SqlSession.Where("id=?", id).First(user).Error
// 	if err != nil {
// 		// 数据库查询步骤出错
// 		return nil
// 	}
// 	return user
// }

func Encryption(str string) string {
	return str
}

func Decryption(str string) string {
	return str
}
