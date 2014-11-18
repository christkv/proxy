package mongo

import (
	"bytes"
	"testing"
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
