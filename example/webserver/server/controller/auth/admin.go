package auth

import (
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/router"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"strings"
	"time"
)

type Admin struct {
	controller.Base

	Authenticate func(account, password string) error
	ErrorCount   map[string]int
}

func (s *Admin) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := &model.SysInfo{
		Name:        s.Config.Name,
		BackVersion: s.Config.GetArgs().ModuleVersion().ToString(),
	}
	a.Success(data)
}

func (s *Admin) GetInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取系统信息")
	function.SetNote("获取服务系统当前相关信息")
	function.SetOutputExample(&model.SysInfo{
		Name:         "WEB服务器",
		BackVersion:  "1.0.1.0",
		FrontVersion: "1.0.1.8",
	})
	function.IgnoreToken(true)
	function.SetContentType("")

	catalog := a.CreateCatalog("平台信息", "平台信息相关接口")
	catalog.SetFunction(function)

	return function
}

func (s *Admin) GetCaptcha(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.CaptchaFilter{
		Mode:   base64Captcha.CaptchaModeNumberAlphabet,
		Length: 4,
		Width:  100,
		Height: 30,
	}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	captchaConfig := base64Captcha.ConfigCharacter{
		Mode:               filter.Mode,
		Height:             filter.Height,
		Width:              filter.Width,
		CaptchaLen:         filter.Length,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
	}
	captchaId, captchaValue := base64Captcha.GenerateCaptcha("", captchaConfig)

	data := &model.Captcha{
		ID:       captchaId,
		Value:    base64Captcha.CaptchaWriteToBase64Encoding(captchaValue),
		Required: s.captchaRequired(a.RIP()),
	}
	randKey := a.RandKey()
	if randKey != nil {
		publicKey, err := randKey.PublicKey()
		if err == nil {
			keyVal, err := publicKey.SaveToMemory()
			if err == nil {
				data.RsaKey = string(keyVal)
			}
		}
	}

	a.Success(data)
}

func (s *Admin) GetCaptchaDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取验证码")
	function.SetNote("获取用户登陆需要的验证码信息")
	function.SetInputExample(&model.CaptchaFilter{
		Mode:   base64Captcha.CaptchaModeNumberAlphabet,
		Length: 4,
		Width:  100,
		Height: 30,
	})
	function.SetOutputExample(&model.Captcha{
		ID:     "GKSVhVMRAHsyVuXSrMYs",
		Value:  "data:image/png;base64,iVBOR...",
		RsaKey: "-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----",
	})
	function.IgnoreToken(true)

	catalog := a.CreateCatalog("权限管理", "系统授权相关接口")
	catalog.SetFunction(function)

	return function
}

func (s *Admin) Login(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.LoginFilter{}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputError, err)
		return
	}

	requireCaptcha := s.captchaRequired(a.RIP())
	err = filter.Check(requireCaptcha)
	if err != nil {
		a.Error(errors.InputInvalid, err)
		return
	}

	if requireCaptcha {
		if !base64Captcha.VerifyCaptcha(filter.CaptchaId, filter.CaptchaValue) {
			a.Error(errors.LoginCaptchaInvalid)
			return
		}
	}

	if s.Authenticate == nil {
		a.Error(errors.Exception, "not auth provider")
		return
	}

	pwd := filter.Password
	if strings.ToLower(filter.Encryption) == "rsa" {
		decryptedPwd, err := a.RandKey().DecryptData(filter.Password)
		if err != nil {
			a.Error(errors.LoginPasswordInvalid, err)
			s.increaseErrorCount(a.RIP())
			return
		}
		pwd = string(decryptedPwd)
	}

	err = s.Authenticate(filter.Account, pwd)
	if err != nil {
		a.Error(errors.LoginAccountOrPasswordInvalid, err)
		s.increaseErrorCount(a.RIP())
		return
	}

	now := time.Now()
	token := &model.Token{
		ID:          a.GenerateGuid(),
		UserAccount: filter.Account,
		LoginIP:     a.RIP(),
		LoginTime:   now,
		ActiveTime:  now,
	}
	err = s.DbToken.Set(token)
	if err != nil {
		a.Error(errors.Exception, err)
		return
	}

	login := &model.Login{
		Token:   token.ID,
		Account: token.UserAccount,
	}

	a.Success(login)
	s.clearErrorCount(a.RIP())
}

func (s *Admin) LoginDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("用户登录")
	function.SetNote("通过用户账号及密码进行登录获取凭证")
	function.SetInputExample(&model.LoginFilter{
		Account:      "admin",
		Password:     "1",
		CaptchaId:    "r4kcmz2E12e0qJQOvqRB",
		CaptchaValue: "1e35",
		Encryption:   "",
	})
	function.SetOutputExample(&model.Login{
		Token: "71b9b7e2ac6d4166b18f414942ff3481",
	})
	function.IgnoreToken(true)

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("权限管理", "系统授权相关接口").SetFunction(function)

	return function
}

func (s *Admin) captchaRequired(ip string) bool {
	if s.ErrorCount == nil {
		return false
	}

	count, ok := s.ErrorCount[ip]
	if ok {
		if count < 3 {
			return false
		} else {
			return true
		}
	}

	return false
}

func (s *Admin) increaseErrorCount(ip string) {
	if s.ErrorCount == nil {
		return
	}

	count := 1
	v, ok := s.ErrorCount[ip]
	if ok {
		count += v
	}

	s.ErrorCount[ip] = count
}

func (s *Admin) clearErrorCount(ip string) {
	if s.ErrorCount == nil {
		return
	}

	_, ok := s.ErrorCount[ip]
	if ok {
		delete(s.ErrorCount, ip)
	}
}
