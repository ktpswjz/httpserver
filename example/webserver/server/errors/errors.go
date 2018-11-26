package errors

import "github.com/ktpswjz/httpserver/types"

var (
	Success       = types.NewError(0, "")
	Unknown       = types.NewError(1, "未知错误")
	Exception     = types.NewError(2, "内部异常")
	InternalError = types.NewError(3, "内部错误")
	NotExist      = types.NewError(4, "不存在")
	Exist         = types.NewError(5, "已存在")
	NotSupport    = types.NewError(6, "不支持的操作")

	InputError   = types.NewError(11, "参数错误")
	InputInvalid = types.NewError(12, "参数无效")

	LoginCaptchaInvalid           = types.NewError(10001, "验证码无效")
	LoginAccountNotExit           = types.NewError(10002, "账号不存在")
	LoginPasswordInvalid          = types.NewError(10003, "密码不正确")
	LoginAccountOrPasswordInvalid = types.NewError(10004, "账号或密码不正确")

	AuthNoToken      = types.NewError(20001, "缺少凭证")
	AuthTokenInvalid = types.NewError(20001, "凭证无效")
	AuthTokenIllegal = types.NewError(20001, "凭证非法")
)
