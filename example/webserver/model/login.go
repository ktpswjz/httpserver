package model

type Login struct {
	Token 	string `json:"token" note:"接口访问凭证" example:"7faf10b0bde847c9905c93966594c82b"`
}

type LoginFilter struct {
	Account string `json:"account" required:"true" note:"账号名称"`
	Password string `json:"password" required:"true" note:"账号密码"`
	CaptchaId string `json:"captchaId" required:"true" note:"验证码ID"`
	CaptchaValue string `json:"captchaValue" required:"true" note:"验证码"`
	Encryption string `json:"encryption" note:"密码加密方法: 空-明文(默认); rsa-RSA密文(公钥通过调用获取验证码接口获取)"`
}