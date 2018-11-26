package document

type Function interface {
	GetName() string
	SetNote(note string)
	GetNote() string
	SetInputExample(example interface{})
	GetInputExample() interface{}
	SetOutputExample(example interface{})
	GetOutputExample() interface{}
	AddInputHeader(name, note string, required bool, values []string)
	AddInputQuery(name, note string, required bool)
	GetInputHeader() []*ModelHeader
	GetInputQuery() []*ModelQuery
	GetOutputHeader() []*ModelHeader
	IgnoreToken(ignore bool)
	IsIgnoreToken() bool
	SetContentType(contentType string)
	GetContentType() string
}

type innerFunction struct {
	name string // 接口名称
	note string // 接口说明

	inputHeaders  []*ModelHeader // 输入请求头部
	inputQueries  []*ModelQuery  // 输入请求参数
	inputExample  interface{}    // 输入参数示例
	outputHeaders []*ModelHeader // 输出头部
	outputExample interface{}    // 输出参数示例

	ignoreToken bool   // 是否忽略凭证
	contentType string // 请求体类型
}

func (s *innerFunction) GetName() string {
	return s.name
}

func (s *innerFunction) SetNote(note string) {
	s.note = note
}

func (s *innerFunction) GetNote() string {
	return s.note
}

func (s *innerFunction) SetInputExample(example interface{}) {
	s.inputExample = example
}

func (s *innerFunction) GetInputExample() interface{} {
	return s.inputExample
}

func (s *innerFunction) SetOutputExample(example interface{}) {
	s.outputExample = example
}

func (s *innerFunction) GetOutputExample() interface{} {
	return s.outputExample
}

func (s *innerFunction) AddInputHeader(name, note string, required bool, values []string) {
	header := &ModelHeader{
		Name:     name,
		Note:     note,
		Required: required,
		Values:   make([]string, 0),
	}
	valueCount := len(values)
	if valueCount > 0 {
		header.DefaultValue = values[0]
		for i := 0; i < valueCount; i++ {
			header.Values = append(header.Values, values[i])
		}
	}

	s.inputHeaders = append(s.inputHeaders, header)
}

func (s *innerFunction) AddInputQuery(name, note string, required bool) {
	query := &ModelQuery{
		Name:     name,
		Note:     note,
		Required: required,
	}

	s.inputQueries = append(s.inputQueries, query)
}

func (s *innerFunction) AddOutputHeader(name, value string) {
	header := &ModelHeader{
		Name:         name,
		DefaultValue: value,
		Values:       make([]string, 0),
	}

	s.outputHeaders = append(s.outputHeaders, header)
}

func (s *innerFunction) GetInputHeader() []*ModelHeader {
	return s.inputHeaders
}

func (s *innerFunction) GetInputQuery() []*ModelQuery {
	return s.inputQueries
}

func (s *innerFunction) GetOutputHeader() []*ModelHeader {
	return s.outputHeaders
}

func (s *innerFunction) IgnoreToken(ignore bool) {
	s.ignoreToken = ignore
}

func (s *innerFunction) IsIgnoreToken() bool {
	return s.ignoreToken
}

func (s *innerFunction) SetContentType(contentType string) {
	s.contentType = contentType
}

func (s *innerFunction) GetContentType() string {
	return s.contentType
}
