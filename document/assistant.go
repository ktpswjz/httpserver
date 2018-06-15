package document

type Assistant interface {
	CreateCatalog(name, note string) Catalog
	CreateFunction(name string) Function
}

type innerAssistant struct {
	catalogs []*innerCatalog
	method string
}

func (s *innerAssistant) CreateCatalog(name, note string) Catalog  {
	catalog := &innerCatalog{}
	catalog.name = name
	catalog.note = note
	s.catalogs = append(s.catalogs, catalog)

	return catalog
}

func (s *innerAssistant) CreateFunction(name string) Function  {
	function := &innerFunction{
		name: name,
		ignoreToken: false,
	}
	if s.method == "POST" {
		function.contentType = "application/json"
	}
	function.inputHeaders = make([]*ModelHeader, 0)
	function.inputQueries = make([]*ModelQuery, 0)
	function.outputHeaders = make([]*ModelHeader, 0)

	function.AddOutputHeader("access-control-allow-origin ", "*")
	function.AddOutputHeader("content-type", "application/json;charset=utf-8")

	return function
}

