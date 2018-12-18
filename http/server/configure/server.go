package configure

type Server struct {
	Http  Http   `json:"http" note:"HTTP服务"`
	Https Https  `json:"https" note:"HTTPS服务"`
	Proxy string `json:"proxy" note:"代理服务器IP地址（客户端不是来自代理服务器时，远程地址为当前连接地址）"`
}
