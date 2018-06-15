package document

type Catalog interface {
	SetFunction(f Function)
	CreateChild(name, note string) Catalog
}

type innerCatalog struct {
	function Function
	catalogs []*innerCatalog
	parent *innerCatalog

	name string
	note string
}

func (s *innerCatalog) SetFunction(f Function) {
	s.function = f
}

func (s *innerCatalog) CreateChild(name, note string) Catalog {
	catalog := &innerCatalog{}
	catalog.name = name
	catalog.note = note
	catalog.parent = s
	s.catalogs = append(s.catalogs, catalog)

	return catalog
}
