package mongo

import (
	"bytes"
	// "strings"
	"testing"
	"time"
)

func TestSimpleEmptyDocumentSerialization(t *testing.T) {
	var expectedBuffer = []byte{5, 0, 0, 0, 0}
	document := NewDocument()
	bson, err := Serialize(document)

	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document")
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	t.Logf("[%v]", obj)

	if obj.FieldCount() != 0 {
		t.Errorf("Failed to deserialize the map")
	}
}

func TestSimpleInt32Serialization(t *testing.T) {
	var expectedBuffer = []byte{14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	document := NewDocument()
	document.Add("int", int32(10))
	bson, err := Serialize(document)

	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document")
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	validateIntField(t, obj, "int", int32(10))
}

func TestSimpleStringSerialization(t *testing.T) {
	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}
	document := NewDocument()
	document.Add("string", "hello world")
	bson, err := Serialize(document)

	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document")
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	validateStringField(t, obj, "string", "hello world")
}

func TestSimpleStringAndIntSerialization(t *testing.T) {
	var expectedBuffer = []byte{38, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	document := NewDocument()
	document.Add("string", "hello world")
	document.Add("int", int32(10))
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 2)
	validateStringField(t, obj, "string", "hello world")
	validateIntField(t, obj, "int", 10)
}

func TestSimpleNestedDocumentSerialization(t *testing.T) {
	var expectedBuffer = []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}
	document := NewDocument()
	subdocument := NewDocument()
	subdocument.Add("int", int32(10))
	document.Add("string", "hello world")
	document.Add("doc", subdocument)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 2)
	validateStringField(t, obj, "string", "hello world")
	validateIntField(t, subDocument(obj, "doc"), "int", 10)
}

func TestSimpleArraySerialization(t *testing.T) {
	var expectedBuffer = []byte{35, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 23, 0, 0, 0, 2, 48, 0, 2, 0, 0, 0, 97, 0, 2, 49, 0, 2, 0, 0, 0, 98, 0, 0, 0}
	document := NewDocument()
	array := make([]interface{}, 0)
	array = append(array, "a")
	array = append(array, "b")
	document.Add("array", array)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	a, _ := obj.GetArray("array")
	validateString(t, a[0], "a")
	validateString(t, a[1], "b")
}

func TestSimpleBinarySerialization(t *testing.T) {
	var expectedBuffer = []byte{26, 0, 0, 0, 5, 98, 105, 110, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0}
	document := NewDocument()
	bin := &Binary{0, []byte("hello world")}
	document.Add("bin", bin)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	validateBinaryField(t, obj, "bin", bin)
}

func TestMixedDocumentSerialization(t *testing.T) {
	var expectedBuffer = []byte{51, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 39, 0, 0, 0, 5, 48, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 3, 49, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0, 0}
	subdocument := NewDocument()
	subdocument.Add("a", int32(1))
	document := NewDocument()
	array := make([]interface{}, 0)
	bin := &Binary{0, []byte("hello world")}
	array = append(array, bin)
	array = append(array, subdocument)
	document.Add("array", array)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	a, _ := obj.GetArray("array")
	validateBinary(t, a[0], bin)
	validateIntField(t, toDocument(t, a[1]), "a", 1)
}

func TestObjectIdSerialization(t *testing.T) {
	var expectedBuffer = []byte{21, 0, 0, 0, 7, 105, 100, 0, 49, 50, 51, 52, 53, 54, 55, 56, 49, 50, 51, 52, 0}
	document := NewDocument()
	objectid := &ObjectId{[]byte("123456781234")}
	document.Add("id", objectid)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	validateObjectIdField(t, obj, "id", objectid)
}

func TestJavascriptNoScopeSerialization(t *testing.T) {
	var expectedBuffer = []byte{34, 0, 0, 0, 13, 106, 115, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 0}
	document := NewDocument()
	js := &Javascript{"var a = function(){}"}
	document.Add("js", js)

	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	validateJavascriptField(t, obj, "js", js)
}

func TestJavascriptWithScopeSerialization(t *testing.T) {
	var expectedBuffer = []byte{50, 0, 0, 0, 15, 106, 115, 0, 41, 0, 0, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0}
	scope := NewDocument()
	scope.Add("a", int32(1))
	document := NewDocument()
	js := &JavascriptWScope{"var a = function(){}", scope}
	document.Add("js", js)

	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	j := validateJavascriptFieldWScope(t, obj, "js", js)
	validateIntField(t, j, "a", int32(1))
}

func TestMinMaxSerialization(t *testing.T) {
	var expectedBuffer = []byte{15, 0, 0, 0, 255, 109, 105, 110, 0, 127, 109, 97, 120, 0, 0}

	// serializeAndPrint('min and max', {min: new MinKey(), max: new MaxKey()});
	document := NewDocument()
	document.Add("min", &Min{})
	document.Add("max", &Max{})

	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Errorf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 2)
	validateMaxField(t, obj, "max")
	validateMinField(t, obj, "min")
}

func TestDateAndTimeSerialization(t *testing.T) {
	var expectedBuffer = []byte{31, 0, 0, 0, 9, 111, 110, 101, 0, 160, 134, 1, 0, 0, 0, 0, 0, 9, 116, 119, 111, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

	// Actual document
	document := NewDocument()
	document.Add("one", &Date{100000})
	document.Add("two", time.Unix(100000, 0))
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestBufferSerialization(t *testing.T) {
	var expectedBuffer = []byte{24, 0, 0, 0, 5, 98, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0}

	// Actual document
	document := NewDocument()
	document.Add("b", []byte("hello world"))
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestTimestampSerialization(t *testing.T) {
	var expectedBuffer = []byte{16, 0, 0, 0, 17, 116, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

	// Actual document
	document := NewDocument()
	document.Add("t", &Timestamp{100000})
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestInt64AndUInt64Serialization(t *testing.T) {
	var expectedBuffer = []byte{27, 0, 0, 0, 18, 111, 0, 255, 255, 255, 255, 255, 255, 255, 255, 18, 116, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

	// Actual document
	document := NewDocument()
	document.Add("o", int64(-1))
	document.Add("t", uint64(100000))
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestFloat64Serialization(t *testing.T) {
	var expectedBuffer = []byte{16, 0, 0, 0, 1, 111, 0, 31, 133, 235, 81, 184, 30, 9, 64, 0}

	// Actual document
	document := NewDocument()
	document.Add("o", float64(3.14))
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestFloat32Serialization(t *testing.T) {
	var expectedBuffer = []byte{16, 0, 0, 0, 1, 111, 0, 102, 102, 102, 102, 102, 102, 246, 191, 0}

	// Actual document
	document := NewDocument()
	document.Add("o", float32(-1.4))
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestBooleanSerialization(t *testing.T) {
	var expectedBuffer = []byte{13, 0, 0, 0, 8, 111, 0, 1, 8, 116, 0, 0, 0}

	// Actual document
	document := NewDocument()
	document.Add("o", true)
	document.Add("t", false)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestNilSerialization(t *testing.T) {
	var expectedBuffer = []byte{8, 0, 0, 0, 10, 111, 0, 0}

	// Actual document
	document := NewDocument()
	document.Add("o", nil)
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}

func TestRegExpSerialization(t *testing.T) {
	var expectedBuffer = []byte{17, 0, 0, 0, 11, 111, 0, 91, 116, 101, 115, 116, 93, 0, 105, 0, 0}

	// Actual document
	document := NewDocument()
	document.Add("o", &RegExp{"[test]", "i"})
	bson, err := Serialize(document)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Errorf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Errorf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Errorf("Illegal BSON returned")
	}
}
