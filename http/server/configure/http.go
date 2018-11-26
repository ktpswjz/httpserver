package configure

type Http struct {
	Enabled         bool   `json:"enabled" note:"是否启用"`
	BehindProxy     bool   `json:"behindProxy" note:"是否位于代理服务器之后"`
	Address         string `json:"address" note:"监听地址，空表示监听所有地址"`
	Port            string `json:"port" note:"监听端口号"`
	RedirectToHttps bool   `json:"redirectToHttps" note:"是否重定向到https"`
}
