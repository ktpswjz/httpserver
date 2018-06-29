package configure

type Certificate struct {
	File     string `json:"file" note:"证书路径年"`
	Password string `json:"password" note:"证书秘密"`
}
