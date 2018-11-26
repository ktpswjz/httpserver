package document

type ModelFunction struct {
	ID            string         `json:"id"`       // 接口标识
	Name          string         `json:"name"`     // 接口名称
	Note          string         `json:"note"`     // 接口说明
	Method        string         `json:"method"`   // 接口方法
	Path          string         `json:"path"`     // 接口地址
	FullPath      string         `json:"fullPath"` // 接口地址
	InputHeaders  []*ModelHeader `json:"inputHeaders"`
	InputQueries  []*ModelQuery  `json:"inputQueries"`
	InputModel    *ModelArgument `json:"inputModel"`
	InputSample   interface{}    `json:"inputSample"`
	OutputHeaders []*ModelHeader `json:"outputHeaders"`
	OutputModel   *ModelArgument `json:"outputModel"`
	OutputSample  interface{}    `json:"outputSample"`
}
