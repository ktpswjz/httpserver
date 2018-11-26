package document

type ModelHeader struct {
	Name         string   `json:"name"`     // 名称
	Note         string   `json:"note"`     // 说明
	Required     bool     `json:"required"` // 必填
	Values       []string `json:"values"`   // 有效值
	DefaultValue string   `json:"defaultValue"`
}
