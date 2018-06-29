package configure

type Server struct {
	Http 	Http 	`json:"http" note:"HTTP服务"`
	Https	Https	`json:"https" note:"HTTPS服务"`
}