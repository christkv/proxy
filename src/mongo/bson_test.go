package mongo

import (
	"bytes"
	"testing"
	"time"
)

func TestSimpleEmptyDocumentSerialization(t *testing.T) {
	var expectedBuffer = []byte{5, 0, 0, 0, 0}
	document := make(map[string]interface{})
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
}

func TestSimpleInt32Serialization(t *testing.T) {
	var expectedBuffer = []byte{14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	document := make(map[string]interface{})
	document["int"] = int32(10)
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
}

func TestSimpleStringSerialization(t *testing.T) {
	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}
	document := make(map[string]interface{})
	document["string"] = "hello world"
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
}

func TestSimpleStringAndIntSerialization(t *testing.T) {
	var expectedBuffer = []byte{38, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	document := make(map[string]interface{})
	document["string"] = "hello world"
	document["int"] = int32(10)
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

func TestSimpleNestedDocumentSerialization(t *testing.T) {
	var expectedBuffer = []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}
	document := make(map[string]interface{})
	subdocument := make(map[string]interface{})
	subdocument["int"] = int32(10)
	document["string"] = "hello world"
	document["doc"] = subdocument
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

func TestSimpleArraySerialization(t *testing.T) {
	var expectedBuffer = []byte{35, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 23, 0, 0, 0, 2, 48, 0, 2, 0, 0, 0, 97, 0, 2, 49, 0, 2, 0, 0, 0, 98, 0, 0, 0}
	document := make(map[string]interface{})
	array := make([]interface{}, 0)
	array = append(array, "a")
	array = append(array, "b")
	document["array"] = array
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

func TestSimpleBinarySerialization(t *testing.T) {
	var expectedBuffer = []byte{26, 0, 0, 0, 5, 98, 105, 110, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0}
	document := make(map[string]interface{})
	document["bin"] = &Binary{0, []byte("hello world")}
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

func TestMixedDocumentSerialization(t *testing.T) {
	var expectedBuffer = []byte{51, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 39, 0, 0, 0, 5, 48, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 3, 49, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0, 0}
	subdocument := make(map[string]interface{})
	subdocument["a"] = int32(1)
	document := make(map[string]interface{})
	array := make([]interface{}, 0)
	array = append(array, &Binary{0, []byte("hello world")})
	array = append(array, subdocument)
	document["array"] = array

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

func TestObjectIdSerialization(t *testing.T) {
	var expectedBuffer = []byte{21, 0, 0, 0, 7, 105, 100, 0, 49, 50, 51, 52, 53, 54, 55, 56, 49, 50, 51, 52, 0}
	document := make(map[string]interface{})
	document["id"] = &ObjectId{[]byte("123456781234")}

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

func TestJavascriptNoScopeSerialization(t *testing.T) {
	var expectedBuffer = []byte{34, 0, 0, 0, 13, 106, 115, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 0}
	document := make(map[string]interface{})
	document["js"] = &Javascript{"var a = function(){}"}

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

func TestJavascriptWithScopeSerialization(t *testing.T) {
	var expectedBuffer = []byte{50, 0, 0, 0, 15, 106, 115, 0, 41, 0, 0, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0}

	scope := make(map[string]interface{})
	scope["a"] = int32(1)
	document := make(map[string]interface{})
	document["js"] = &JavascriptWScope{"var a = function(){}", scope}

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

func TestMinMaxSerialization(t *testing.T) {
	var expectedBuffer = []byte{50, 0, 0, 0, 15, 106, 115, 0, 41, 0, 0, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0}

	scope := make(map[string]interface{})
	scope["a"] = int32(1)
	document := make(map[string]interface{})
	document["js"] = &JavascriptWScope{"var a = function(){}", scope}

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

func TestDateAndTimeSerialization(t *testing.T) {
	var expectedBuffer = []byte{31, 0, 0, 0, 9, 111, 110, 101, 0, 160, 134, 1, 0, 0, 0, 0, 0, 9, 116, 119, 111, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

	// Actual document
	document := make(map[string]interface{})
	document["one"] = &Date{100000}
	document["two"] = time.Unix(100000, 0)
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
	document := make(map[string]interface{})
	document["b"] = []byte("hello world")
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
	document := make(map[string]interface{})
	document["t"] = &Timestamp{100000}
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
	document := make(map[string]interface{})
	document["o"] = int64(-1)
	document["t"] = uint64(100000)
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
	document := make(map[string]interface{})
	document["o"] = float64(3.14)
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
	document := make(map[string]interface{})
	document["o"] = float32(-1.4)
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
	document := make(map[string]interface{})
	document["o"] = true
	document["t"] = false
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
	document := make(map[string]interface{})
	document["o"] = nil
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
	document := make(map[string]interface{})
	document["o"] = &RegExp{"[test]", "i"}
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