package document

import "testing"

func TestParse(t *testing.T) {
	example := &Example{}

	model := example.ParseModel(nil)
	if model != nil {
		t.Fatal(model)
	}

	var interfaceModel interface{}
	model = example.ParseModel(&interfaceModel)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "interface", "", "", false, 0)

	interfaceModel = "123"
	model = example.ParseModel(&interfaceModel)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "string", "", "", false, 0)

	model = example.ParseModel(interfaceModel)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "string", "", "", false, 0)

	intModel := 8
	model = example.ParseModel(&intModel)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "int", "", "", false, 0)

	stringModel := "s"
	model = example.ParseModel(&stringModel)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "string", "", "", false, 0)
	model = example.ParseModel(stringModel)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "string", "", "", false, 0)

	argumentBase := ArgumentBase{}
	model = example.ParseModel(&argumentBase)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "ArgumentBase", "", "", false, 1)
	checkModel(t, model.Childs[0], "string", "name", "名称", false, 0)

	argument1 := Argument1{}
	model = example.ParseModel(&argument1)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "Argument1", "", "", false, 3)
	checkModel(t, model.Childs[0], "string", "name", "名称", false, 0)
	checkModel(t, model.Childs[1], "uint64", "id", "标识ID", true, 0)
	checkModel(t, model.Childs[2], "interface", "age", "", false, 0)

	argument11 := Argument1{Age: 5}
	model = example.ParseModel(&argument11)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "Argument1", "", "", false, 3)
	checkModel(t, model.Childs[2], "int", "age", "", false, 0)

	argument2 := Argument2{}
	model = example.ParseModel(&argument2)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "Argument2", "", "", false, 2)
	checkModel(t, model.Childs[0], "uint64", "id", "标识ID2", true, 0)
	checkModel(t, model.Childs[1], "ArgumentBase", "base", "", false, 1)
	checkModel(t, model.Childs[1].Childs[0], "string", "name", "名称", false, 0)

	argument4 := Argument4{}
	model = example.ParseModel(&argument4)
	if model == nil {
		t.Fatal("unkown")
	}
	if model.Type != "Argument4" {
		t.Error("argument4 type: expect=", "Argument4", ", actual=", model.Type)
	}
	if len(model.Childs) != 1 {
		t.Error("argument4 count: expect=", 1, ", actual=", len(model.Childs))
	}
	modelChild := model.Childs[0]
	if modelChild.Name != "id" {
		t.Error("argument4-0 name: expect=", "id", ", actual=", modelChild.Name)
	}
	if modelChild.Type != "*uint64" {
		t.Error("argument4-0 type: expect=", "*uint64", ", actual=", modelChild.Type)
	}
}

func TestParseOfSlice(t *testing.T) {
	example := &Example{}
	model := example.ParseModel(nil)
	if model != nil {
		t.Fatal(model)
	}

	var interfaceSlice []interface{}
	model = example.ParseModel(&interfaceSlice)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "interface[]", "", "", false, 0)

	var intSlice []int
	model = example.ParseModel(&intSlice)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "int[]", "", "", false, 0)

	var stringSlice []string
	model = example.ParseModel(&stringSlice)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "string[]", "", "", false, 0)
	var stringSlice2 []*string
	model = example.ParseModel(stringSlice2)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "string[]", "", "", false, 0)

	var structSlice []ArgumentBase
	model = example.ParseModel(structSlice)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "ArgumentBase[]", "", "", false, 1)
	checkModel(t, model.Childs[0], "ArgumentBase", "", "", false, 1)
	checkModel(t, model.Childs[0].Childs[0], "string", "name", "名称", false, 0)

	var structSlice2 []*ArgumentBase
	model = example.ParseModel(structSlice2)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "ArgumentBase[]", "", "", false, 1)
	checkModel(t, model.Childs[0], "ArgumentBase", "", "", false, 1)
	checkModel(t, model.Childs[0].Childs[0], "string", "name", "名称", false, 0)

	var structSlice3 []*ArgumentBase
	model = example.ParseModel(&structSlice3)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "ArgumentBase[]", "", "", false, 1)
	checkModel(t, model.Childs[0], "ArgumentBase", "", "", false, 1)
	checkModel(t, model.Childs[0].Childs[0], "string", "name", "名称", false, 0)

	var structSlice4 []ArgumentBase
	model = example.ParseModel(&structSlice4)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "ArgumentBase[]", "", "", false, 1)
	checkModel(t, model.Childs[0], "ArgumentBase", "", "", false, 1)
	checkModel(t, model.Childs[0].Childs[0], "string", "name", "名称", false, 0)

	var argument3 Argument3
	model = example.ParseModel(&argument3)
	if model == nil {
		t.Fatal("unkown")
	}
	checkModel(t, model, "Argument3", "", "", false, 5)
	checkModel(t, model.Childs[0], "uint64", "id", "标识ID2", true, 0)
	checkModel(t, model.Childs[1], "ArgumentBase", "base", "", false, 1)
	checkModel(t, model.Childs[1].Childs[0], "string", "name", "名称", false, 0)
	checkModel(t, model.Childs[2], "uint64", "id3", "标识ID3", false, 0)
	checkModel(t, model.Childs[3], "string[]", "addresses", "地址", false, 0)
	checkModel(t, model.Childs[4], "ArgumentBase[]", "names", "", false, 1)
	checkModel(t, model.Childs[4].Childs[0], "ArgumentBase", "", "", false, 1)
	checkModel(t, model.Childs[4].Childs[0].Childs[0], "string", "name", "名称", false, 0)
}

func checkModel(t *testing.T, m *ModelArgument, argType, argName, argNote string, argRequired bool, argChildCount int) {
	if m.Type != argType {
		t.Error("type: expect=", argType, ", actual=", m.Type)
	}
	if m.Name != argName {
		t.Error("name: expect=", argName, ", actual=", m.Name)
	}
	if m.Note != argNote {
		t.Error("note: expect=", argNote, ", actual=", m.Note)
	}
	if m.Required != argRequired {
		t.Error("required: expect=", argRequired, ", actual=", m.Required)
	}
	if len(m.Childs) != argChildCount {
		t.Fatal("childs: expect=", argChildCount, ", actual=", len(m.Childs))
	}
}

type ArgumentBase struct {
	Name string `json:"name" note:"名称"`
}

type Argument1 struct {
	ArgumentBase

	ID  uint64      `json:"id" required:"true" note:"标识ID"`
	Age interface{} `json:"age"`
}

type Argument2 struct {
	ID uint64 `json:"id" required:"true" note:"标识ID2"`

	Base ArgumentBase `json:"base"`
}

type Argument3 struct {
	Argument2

	ID3       uint64          `json:"id3" required:"false" note:"标识ID3"`
	Addresses []string        `json:"addresses" note:"地址"`
	Names     []*ArgumentBase `json:"names"`
}

type Argument4 struct {
	ID *uint64 `json:"id,omitempty" required:"true" note:"标识"`
}
