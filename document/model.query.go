package document

type ModelQuery struct {
	Name string `json:"name"` 		// 名称
	Note string `json:"note"` 		// 说明
	Required bool `json:"required"`	// 必填
	DefaultValue string `json:"defaultValue"`
}