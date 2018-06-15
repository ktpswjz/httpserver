package configure

type Http struct {
	Enabled bool   `json:"enabled"` // true or false
	Address string `json:"address"` // listen or server address
	Port    string `json:"port"`    // listen or server port number
}