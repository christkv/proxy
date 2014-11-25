package mongo

import (
	"bytes"
	"testing"
	"time"
)

func validateIntField(t *testing.T, obj *Document, name string, expected int32) {
	value, err := obj.FieldAsInt32(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as int32", name)
		return
	}

	if value != expected {
		t.Fatalf("Failed int32 comparison [%v] != [%v]", value, expected)
	}
}

func validateTimeField(t *testing.T, obj *Document, name string, expected time.Time) {
	value, err := obj.FieldAsTime(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as time", name)
		return
	}

	if value.Unix() != expected.Unix() {
		t.Fatalf("Failed time comparison [%v] != [%v]", value, expected)
	}
}

func validateStringField(t *testing.T, obj *Document, name string, expected string) {
	value, err := obj.FieldAsString(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as string", name)
		return
	}

	if value != expected {
		t.Fatalf("Failed string comparison [%v] != [%v]", value, expected)
	}
}

func validateBinaryField(t *testing.T, obj *Document, name string, expected *Binary) {
	value, err := obj.FieldAsBinary(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as binary", name)
		return
	}

	if value.SubType != expected.SubType || bytes.Compare(value.Data, expected.Data) != 0 {
		t.Fatalf("Failed binary comparison [%v] != [%v]", value, expected)
	}
}

func validateBufferField(t *testing.T, obj *Document, name string, expected []byte) {
	value, err := obj.FieldAsBinary(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as binary", name)
		return
	}

	if value.SubType != 0 || bytes.Compare(value.Data, expected) != 0 {
		t.Fatalf("Failed binary comparison [%v] != [%v]", value, expected)
	}
}

func validateTimestampField(t *testing.T, obj *Document, name string, expected *Timestamp) {
	value, err := obj.FieldAsTimestamp(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as timestamp", name)
		return
	}

	if expected.Value != value.Value {
		t.Fatalf("Failed timestamp comparison [%v] != [%v]", value, expected)
	}
}

func validateObjectIdField(t *testing.T, obj *Document, name string, expected *ObjectId) {
	value, err := obj.FieldAsObjectId(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as objectid", name)
		return
	}

	if bytes.Compare(value.Id, value.Id) != 0 {
		t.Fatalf("Failed objectid comparison [%v] != [%v]", value, expected)
	}
}

func validateInt64Field(t *testing.T, obj *Document, name string, expected int64) {
	value, err := obj.FieldAsInt64(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as int64", name)
		return
	}

	if expected != value {
		t.Fatalf("Failed int64 comparison [%v] != [%v]", value, expected)
	}
}

func validateUInt64Field(t *testing.T, obj *Document, name string, expected uint64) {
	value, err := obj.FieldAsUInt64(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as uint64", name)
		return
	}

	if expected != value {
		t.Fatalf("Failed uint64 comparison [%v] != [%v]", value, expected)
	}
}

func validateFloat64Field(t *testing.T, obj *Document, name string, expected float64) {
	value, err := obj.FieldAsFloat64(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as float64", name)
		return
	}

	if expected != value {
		t.Fatalf("Failed float64 comparison [%v] != [%v]", value, expected)
	}
}

func validateFloat32Field(t *testing.T, obj *Document, name string, expected float32) {
	value, err := obj.FieldAsFloat32(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as float32", name)
		return
	}

	if expected != value {
		t.Fatalf("Failed float32 comparison [%v] != [%v]", value, expected)
	}
}

func validateBooleanField(t *testing.T, obj *Document, name string, expected bool) {
	value, err := obj.FieldAsBool(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as bool", name)
		return
	}

	if expected != value {
		t.Fatalf("Failed bool comparison [%v] != [%v]", value, expected)
	}
}

func validateNilField(t *testing.T, obj *Document, name string) {
	value, err := obj.FieldAsNil(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as bool", name)
		return
	}

	if value != nil {
		t.Fatalf("Failed nil comparison [%v]", value)
	}
}

func validateRegExpField(t *testing.T, obj *Document, name string, expected *RegExp) {
	value, err := obj.FieldAsRegExp(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as RegExp", name)
		return
	}

	if value.Options != expected.Options || value.Pattern != expected.Pattern {
		t.Fatalf("Failed RegExp comparison [%v] != [%v]", value, expected)
	}
}

func validateJavascriptField(t *testing.T, obj *Document, name string, expected *Javascript) {
	value, err := obj.FieldAsJavascript(name)
	if err != nil {
		t.Fatalf("Failed to retrieve value [%v] as JavascriptField", name)
		return
	}

	if value.Code != expected.Code {
		t.Fatalf("Failed JavascriptField comparison [%v] != [%v]", value, expected)
	}
}

func validateJavascriptFieldWScope(t *testing.T, obj *Document, name string, expected *JavascriptWScope) *Document {
	value, err := obj.FieldAsJavascriptWScope(name)
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
	value, err := obj.FieldAsDocument(name)
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
