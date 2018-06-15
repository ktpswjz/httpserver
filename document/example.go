package document

import (
	"reflect"
	"fmt"
)

type Example struct {

}

func (s *Example) ParseModel(example interface{}) *ModelArgument {
	if example == nil {
		return nil
	}

	model := &ModelArgument{Childs: make([]*ModelArgument, 0)}

	exampleType := reflect.TypeOf(example)
	exampleTypeKind := exampleType.Kind()
	switch exampleTypeKind {
	case reflect.Ptr : {
		s.parseModel(reflect.ValueOf(example).Elem(), model)
		break
	}
	case reflect.Interface,
		reflect.Struct,
		reflect.Array,
		reflect.Slice,
		reflect.Bool,
		reflect.String,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64 : {
		s.parseModel(reflect.ValueOf(example), model)
		break
	}
	default:
		return nil
	}

	return model
}

func (s *Example) parseModel(v reflect.Value, argument *ModelArgument)  {
	if argument == nil {
		return
	}
	if v.Kind() == reflect.Invalid {
		return
	}

	t := v.Type()
	k := t.Kind()
	switch k {
	case reflect.Ptr : {
		s.parseModel(v.Elem(), argument)
		break
	}
	case reflect.Interface: {
		argument.Type = k.String()
		if v.CanInterface() {
			value := reflect.ValueOf(v.Interface())
			if value.Kind() != reflect.Invalid {
				s.parseModel(value, argument)
			}
		}
		break
	}
	case reflect.Struct : {
		if argument.Type == "" {
			argument.Type = t.Name()
		}

		n := v.NumField()
		for i := 0; i < n; i++ {
			valueField := v.Field(i)
			if !valueField.CanInterface() {
				continue
			}

			typeField := t.Field(i)
			if typeField.Anonymous {
				if valueField.CanAddr() {
					s.parseModel(valueField.Addr().Elem(), argument)
				}
			} else {
				child := &ModelArgument{Childs: make([]*ModelArgument, 0)}
				child.Name = typeField.Tag.Get("json")
				if child.Name == "" {
					child.Name = typeField.Name
				}
				child.Type = valueField.Kind().String()
				if typeField.Tag.Get("required") == "true" {
					child.Required = true
				}
				child.Note = typeField.Tag.Get("note")
				argument.Childs = append(argument.Childs, child)

				value := reflect.ValueOf(valueField.Interface())
				if value.Kind() != reflect.Invalid {
					child.Type = value.Type().Name()
					s.parseModel(value, child)
				}
			}
		}
		break
	}
	case reflect.Array: {
		break
	}
	case reflect.Slice: {
		st := t.Elem()
		stk := st.Kind()

		var ste reflect.Type = nil
		if stk == reflect.Ptr {
			ste = st.Elem()
		} else {
			ste = st
		}
		if ste != nil {
			argument.Type =  fmt.Sprintf("%s[]", ste.Name())

			if ste.Kind() == reflect.Struct {
				stet := reflect.New(ste)
				child := &ModelArgument{Childs: make([]*ModelArgument, 0)}
				child.Type = ste.Name()
				argument.Childs = append(argument.Childs, child)
				s.parseModel(stet.Elem(), child)
			}
		} else {
			argument.Type = fmt.Sprintf("%s[]", stk.String())
		}

		if argument.Type == "[]" {
			argument.Type =  fmt.Sprintf("%s[]", stk.String())
		}

		break
	}
	case reflect.Bool,
		reflect.String,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64 : {
		argument.Type = t.Name()
		break
	}
	default: {
		return
	}
	}

}