package configure

type Https struct {
	Enabled     bool        `json:"enabled" note:"是否启用"`            // true or false
	BehindProxy bool        `json:"behindProxy" note:"是否位于代理服务器之后"` // true or false
	Address     string      `json:"address" note:"监听地址，空表示监听所有地址"`  // listen or server address
	Port        string      `json:"port" note:"监听端口号"`              // listen or server port number
	Cert        Certificate `json:"cert" note:"证书"`
}
