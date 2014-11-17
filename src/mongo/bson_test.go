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
