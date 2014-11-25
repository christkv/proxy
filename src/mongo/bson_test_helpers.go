package mongo

import (
	"bytes"
	"testing"
)

func validateIntField(t *testing.T, obj *Document, name string, expected int32) {
	value, err := obj.GetInt32(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as int32", name)
		return
	}

	if value != expected {
		t.Fatalf("Failed int32 comparison [%v] != [%v]", value, expected)
	}
}

func validateStringField(t *testing.T, obj *Document, name string, expected string) {
	value, err := obj.GetString(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as string", name)
		return
	}

	if value != expected {
		t.Fatalf("Failed string comparison [%v] != [%v]", value, expected)
	}
}

func validateBinaryField(t *testing.T, obj *Document, name string, expected *Binary) {
	value, err := obj.GetBinary(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as binary", name)
		return
	}

	if value.SubType != expected.SubType || bytes.Compare(value.Data, expected.Data) != 0 {
		t.Fatalf("Failed binary comparison [%v] != [%v]", value, expected)
	}
}

func validateObjectIdField(t *testing.T, obj *Document, name string, expected *ObjectId) {
	value, err := obj.GetObjectId(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as objectid", name)
		return
	}

	if bytes.Compare(value.Id, value.Id) != 0 {
		t.Fatalf("Failed objectid comparison [%v] != [%v]", value, expected)
	}
}

func validateJavascriptField(t *testing.T, obj *Document, name string, expected *Javascript) {
	value, err := obj.GetJavascriptField(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as JavascriptField", name)
		return
	}

	if value.Code != expected.Code {
		t.Fatalf("Failed JavascriptField comparison [%v] != [%v]", value, expected)
	}
}

func validateJavascriptFieldWScope(t *testing.T, obj *Document, name string, expected *JavascriptWScope) *Document {
	value, err := obj.GetJavascriptWScopeField(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as JavascriptField", name)
		return nil
	}

	if value.Code != expected.Code {
		t.Fatalf("Failed JavascriptField comparison [%v] != [%v]", value, expected)
	}

	return value.Scope
}

func validateMaxField(t *testing.T, obj *Document, name string) {
	_, err := obj.FieldAsMax(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as Max", name)
		return
	}
}

func validateMinField(t *testing.T, obj *Document, name string) {
	_, err := obj.FieldAsMin(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as Max", name)
		return
	}
}

func validateObjectSize(t *testing.T, obj *Document, size int) {
	if obj.FieldCount() != size {
		t.Fatalf("Failed to deserialize the map")
	}
}

func subDocument(obj *Document, name string) *Document {
	value, err := obj.GetDocument(name)
	if err != nil {
		return nil
	}

	return value
}

func toDocument(t *testing.T, val interface{}) *Document {
	switch elem := val.(type) {
	default:
		t.Fatalf("type of value passed into toDocument is not a *Document")
		return nil
	case *Document:
		return elem
	}
}

func validateString(t *testing.T, v1 interface{}, value string) {
	switch elem := v1.(type) {
	default:
		t.Fatalf("type of value passed into validateString is not string")
	case string:
		if elem != value {
			t.Fatalf("Failed string comparison [%v] != [%v]", elem, value)
		}
	}
}

func validateBinary(t *testing.T, v1 interface{}, bin *Binary) {
	switch elem := v1.(type) {
	default:
		t.Fatalf("type of value passed into validateBinary is not a binary")
	case *Binary:
		if elem.SubType != bin.SubType || bytes.Compare(elem.Data, bin.Data) != 0 {
			t.Fatalf("Failed binary comparison [%v] != [%v]", elem, bin)
		}
	}
}

// func validateBinary(t *testing.T, v1 interface{}, bin *Binary) {
//  switch elem := v1.(type) {
//  default:
//    t.Fatalf("type of value passed into validateBinary is not a binary")
//  case *Binary:
//    if elem.SubType != bin.SubType || bytes.Compare(elem.Data, bin.Data) != 0 {
//      t.Fatalf("Failed binary comparison [%v] != [%v]", elem, bin)
//    }
//  }
// }

// func subArray(obj map[string]interface{}, name string) []interface{} {
//  switch elem := obj[name].(type) {
//  default:
//    return nil
//  case []interface{}:
//    return elem
//  }
// }
