package xerr

import (
	"net/http"
)

// 错误码设计
// 第一位表示错误级别, 1 为系统错误, 2 为普通错误
// 第二三位表示服务模块代码
// 第四五位表示具体错误代码

var (
	OK = &CodeError{errCode: 0, errMsg: "OK", httpCode: http.StatusOK}

	// ErrUnKnown        未知错误
	// ErrInternalServer 类型转换，使用标准库函数，使用第三方包等时出现的错误
	// ErrDatabase       使用mysql数据库或者redis数据库出错
	ErrUnKnown        = &CodeError{errCode: 10000, errMsg: "未知错误", httpCode: http.StatusInternalServerError}
	ErrInternalServer = &CodeError{errCode: 10001, errMsg: "内部服务器错误", httpCode: http.StatusInternalServerError}
	ErrDatabase       = &CodeError{errCode: 10002, errMsg: "数据库错误", httpCode: http.StatusInternalServerError}

	// 模块通用错误
	ErrBadRequest      = &CodeError{errCode: 20001, errMsg: "请求参数不合法", httpCode: http.StatusOK}
	ErrGenToken        = &CodeError{errCode: 20002, errMsg: "生成 token 失败", httpCode: http.StatusOK}
	ErrTokenNotFound   = &CodeError{errCode: 20003, errMsg: "用户 token 不存在或已过期", httpCode: http.StatusOK}
	ErrTokenValidation = &CodeError{errCode: 20004, errMsg: "用户 token 无效", httpCode: http.StatusOK}

	// User模块错误
	ErrUserNotFound       = &CodeError{errCode: 20101, errMsg: "用户不存在", httpCode: http.StatusOK}
	ErrUserExist          = &CodeError{errCode: 20102, errMsg: "用户已存在", httpCode: http.StatusOK}
	ErrPasswordIncorrect  = &CodeError{errCode: 20103, errMsg: "密码错误", httpCode: http.StatusOK}
	ErrUsernameValidation = &CodeError{errCode: 20104, errMsg: "用户名不合法", httpCode: http.StatusOK}
	ErrPasswordValidation = &CodeError{errCode: 20105, errMsg: "密码不合法", httpCode: http.StatusOK}

	// Video模块错误
	ErrInvaildFile = &CodeError{errCode: 20201, errMsg: "无效的文件，请检查文件", httpCode: http.StatusOK}
)
