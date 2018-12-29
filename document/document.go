package document

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/ktpswjz/httpserver/types"
	"strings"
)

type Document interface {
	AddFunction(method, path string, handle Handle)
	GetFunction(id string) *ModelFunction
	GetCatalogs(keywords string) []*ModelCatalog
	GetCatalogTree(keywords string) []*ModelCatalogTree
	GenerateCatalogTree()
}

func NewDocument(enable bool, log types.Log) Document {
	instance := &innerDocument{enable: enable}
	instance.functions = make(map[string]*ModelFunction)
	instance.catalogs = make([]*ModelCatalog, 0)
	instance.catalogTree = make([]*ModelCatalogTree, 0)
	instance.SetLog(log)

	return instance
}

type innerDocument struct {
	types.Base
	enable      bool
	functions   map[string]*ModelFunction
	catalogs    []*ModelCatalog
	catalogTree []*ModelCatalogTree
}

func (s *innerDocument) AddFunction(method, path string, handle Handle) {
	if !s.enable {
		return
	}

	if handle == nil {
		return
	}

	assistant := &innerAssistant{method: method}
	assistant.catalogs = make([]*innerCatalog, 0)
	fun := handle(assistant)
	if fun == nil {
		return
	}

	s.appendFunction(fun, method, path, assistant)
}

func (s *innerDocument) GetFunction(id string) *ModelFunction {
	return s.functions[id]
}

func (s *innerDocument) GetCatalogs(keywords string) []*ModelCatalog {
	return s.catalogs
}

func (s *innerDocument) GetCatalogTree(keywords string) []*ModelCatalogTree {
	return s.catalogTree
}

func (s *innerDocument) GenerateCatalogTree() {
	s.catalogTree = make([]*ModelCatalogTree, 0)

	cn := len(s.catalogs)
	for ci := 0; ci < cn; ci++ {
		c := s.catalogs[ci]
		tc := &ModelCatalogTree{Name: c.Name, Note: c.Note, Type: 0}
		tc.Children = make([]*ModelCatalogTree, 0)
		s.catalogTree = append(s.catalogTree, tc)

		s.generateCatalogTree(tc, c)
	}
}

func (s *innerDocument) generateCatalogTree(tree *ModelCatalogTree, catalog *ModelCatalog) {
	cn := len(catalog.Catalogs)
	for ci := 0; ci < cn; ci++ {
		c := catalog.Catalogs[ci]
		tc := &ModelCatalogTree{Name: c.Name, Note: c.Note, Type: 0}
		tc.Children = make([]*ModelCatalogTree, 0)
		tree.Children = append(tree.Children, tc)

		s.generateCatalogTree(tc, c)
	}

	fc := len(catalog.Functions)
	for fi := 0; fi < fc; fi++ {
		f := catalog.Functions[fi]
		tf := &ModelCatalogTree{Name: f.Name, ID: f.ID, Type: 1}
		tfi := s.functions[f.ID]
		if tfi != nil {
			keywords := &strings.Builder{}
			keywords.WriteString(tfi.Method)
			keywords.WriteString(tfi.Path)
			keywords.WriteString(tfi.Note)
			tf.Keywords = keywords.String()
		}

		tree.Children = append(tree.Children, tf)
	}
}

func (s *innerDocument) appendFunction(f Function, method, path string, assistant *innerAssistant) {
	if f == nil {
		return
	}

	example := Example{}
	modelFunction := &ModelFunction{Method: method}
	modelFunction.Path = path
	modelFunction.Name = f.GetName()
	modelFunction.Note = f.GetNote()
	modelFunction.WebSocket = f.IsWebSocket()
	modelFunction.ID = s.generateFunctionId(method, path)
	modelFunction.InputHeaders = make([]*ModelHeader, 0)
	if !f.IsIgnoreToken() {
		tokenHeader := &ModelHeader{
			Name:     "token",
			Note:     "凭证",
			Required: true,
			Values:   make([]string, 0),
		}
		modelFunction.InputHeaders = append(modelFunction.InputHeaders, tokenHeader)
	}
	if f.GetContentType() != "" {
		contentTypeHeader := &ModelHeader{
			Name:         "Content-Type",
			Note:         "内容类型",
			Required:     true,
			Values:       make([]string, 0),
			DefaultValue: f.GetContentType(),
		}
		contentTypeHeader.Values = append(contentTypeHeader.Values, f.GetContentType())
		modelFunction.InputHeaders = append(modelFunction.InputHeaders, contentTypeHeader)
	}

	headers := f.GetInputHeader()
	headerLength := len(headers)
	for i := 0; i < headerLength; i++ {
		modelFunction.InputHeaders = append(modelFunction.InputHeaders, headers[i])
	}
	modelFunction.InputQueries = f.GetInputQuery()
	modelFunction.InputSample = f.GetInputExample()
	modelFunction.InputModel = example.ParseModel(f.GetInputExample())

	modelFunction.OutputHeaders = f.GetOutputHeader()
	output := &types.Result{
		Code:   0,
		Serial: 201805161315480008,
		Error: types.ResultError{
			Summary: "",
			Detail:  "",
		},
		Data: f.GetOutputExample(),
	}
	modelFunction.OutputSample = output
	modelFunction.OutputModel = example.ParseModel(output)

	s.functions[modelFunction.ID] = modelFunction
	s.appendCatalogs(modelFunction.ID, assistant.catalogs)

	s.LogDebug("api(id=", modelFunction.ID, ", path=", path, ") has been ready")
}

func (s *innerDocument) appendCatalogs(functionId string, catalogs []*innerCatalog) {
	n := len(catalogs)
	if n < 1 {
		return
	}

	for i := 0; i < n; i++ {
		f := catalogs[i].function
		if f != nil {
			s.appendCatalog(catalogs[i], functionId, f.GetName())
		}

		s.appendCatalogs(functionId, catalogs[i].catalogs)
	}
}

func (s *innerDocument) appendCatalog(catalog *innerCatalog, functionId, functionName string) {
	if catalog == nil {
		return
	}
	paths := make([]*innerCatalog, 0)
	paths = append(paths, catalog)
	parent := catalog.parent
	for parent != nil {
		paths = append(paths, parent)
		parent = parent.parent
	}
	n := len(paths)

	var currentModelCatalog *ModelCatalog = nil
	var catalogContainer ModelCatalogContainer = s
	for i := n - 1; i >= 0; i-- {
		p := paths[i]
		modelCatalog := catalogContainer.GetChildByName(p.name)
		if modelCatalog == nil {
			currentModelCatalog = catalogContainer.AddChild(p.name, p.note)
		} else {
			currentModelCatalog = modelCatalog
		}
		catalogContainer = currentModelCatalog
	}

	if currentModelCatalog != nil {
		currentModelCatalog.Functions = append(currentModelCatalog.Functions, &ModelCatalogFunction{
			ID:   functionId,
			Name: functionName,
		})
	}
}

func (s *innerDocument) GetChildByName(name string) *ModelCatalog {
	n := len(s.catalogs)
	for i := 0; i < n; i++ {
		if s.catalogs[i].Name == name {
			return s.catalogs[i]
		}
	}

	return nil
}

func (s *innerDocument) AddChild(name, note string) *ModelCatalog {
	catalog := &ModelCatalog{Name: name, Note: note}
	catalog.Catalogs = make([]*ModelCatalog, 0)
	catalog.Functions = make([]*ModelCatalogFunction, 0)

	s.catalogs = append(s.catalogs, catalog)

	return catalog
}

func (s *innerDocument) jsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *innerDocument) generateFunctionId(method, path string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(method + path))
	if err != nil {
		return path
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
