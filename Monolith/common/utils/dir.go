package utils

import (
	"os"

	"go.uber.org/zap"
)

//判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建文件夹
func CreateDir(path string) {
	exist, err := HasDir(path)
	if err != nil {
		zap.L().Fatal("check videos dir failed")
		return
	}

	if !exist {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			zap.L().Fatal("mkdir videos failed")
			return
		}
		zap.L().Info("create videos success")
	} else {
		zap.L().Info("videos dir already exists.")
	}
}
