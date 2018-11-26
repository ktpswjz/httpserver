package document

type ModelArgument struct {
	parent *ModelArgument

	Name     string `json:"name"`     // 名称
	Type     string `json:"type"`     // 类型
	Note     string `json:"note"`     // 说明
	Required bool   `json:"required"` // 必填

	Childs []*ModelArgument `json:"childs"`
}

func (s *ModelArgument) ParentType() string {
	if s.parent == nil {
		return ""
	}

	return s.parent.Type
}
