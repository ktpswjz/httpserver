package document

type ModelCatalogContainer interface {
	GetChildByName(name string) *ModelCatalog
	AddChild(name, note string) *ModelCatalog
}

type ModelCatalog struct {
	Name string `json:"name"`	// 接口名称
	Note string `json:"note"`	// 接口说明

	Catalogs []*ModelCatalog 			`json:"catalogs"`
	Functions []*ModelCatalogFunction 	`json:"functions"`
}

func (s *ModelCatalog) GetChildByName(name string) *ModelCatalog  {
	n := len(s.Catalogs)
	for i := 0; i < n; i++ {
		if s.Catalogs[i].Name == name {
			return s.Catalogs[i]
		}
	}

	return  nil
}

func (s *ModelCatalog) AddChild(name, note string) *ModelCatalog {
	catalog := &ModelCatalog{Name:name, Note:note}
	catalog.Catalogs = make([]*ModelCatalog, 0)
	catalog.Functions = make([]*ModelCatalogFunction, 0)

	s.Catalogs = append(s.Catalogs, catalog)

	return catalog
}
