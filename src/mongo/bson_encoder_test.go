package mongo

import (
	"bytes"
	// "fmt"
	// "strings"
	"reflect"
	"testing"
	// "time"
)

func SerializeTest(t *testing.T, doc interface{}, expectedBuffer []byte) {
	parser := NewBSON()
	// Serialize the document allowing self allocation of buffer
	b, err := parser.Marshall(doc, nil, 0)

	// t.Logf("%s", b)

	// Ensure the buffers match
	if err != nil || len(b) != len(expectedBuffer) || bytes.Compare(b, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned \nexp: %v:%v\ngot: %v:%v", expectedBuffer, len(expectedBuffer), b, len(b))
	}

	// Serialize into pre-allocated buffer
	b = make([]byte, len(expectedBuffer))
	// Serialize the document
	b, err = parser.Marshall(doc, b, 0)
	// Ensure the buffers match
	if err != nil || len(b) != len(expectedBuffer) || bytes.Compare(b, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned \nexp: %v:%v\ngot: %v:%v", expectedBuffer, len(expectedBuffer), b, len(b))
	}
}

func DeserializeTest(t *testing.T, b []byte, empty interface{}, expected interface{}) {
	parser := NewBSON()
	// Deserialize the data into the type
	err := parser.Unmarshal(b, empty)
	if err != nil {
		t.Errorf("[%v] Failed to unmarshal %v into type %v", err, b, expected)
	}

	// Check if this is a document
	switch doc := empty.(type) {
	case *Document:
		switch doc1 := expected.(type) {
		case *Document:
			if doc1.Equal(doc) == false {
				t.Errorf("failed to unmarshal document correctly 4")
			}
		}
	default:
		if reflect.DeepEqual(empty, expected) == false {
			t.Errorf("failed to unmarshal struct correctly 5")
		}
	}
}

func TestOverflowDefaultBufferSize(t *testing.T) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	// Should correctly create bson from struct
	doc := &T2{"hello world hello world hello world hello world hello world hello world", &T1{10}}
	parser := NewBSON()
	parser.Marshall(doc, nil, 0)

	// Should correctly create bson from Document
	document := NewDocument()
	subdocument := NewDocument()
	subdocument.Add("int", int32(10))
	document.Add("string", "hello world hello world hello world hello world hello world hello world")
	document.Add("doc", subdocument)
	parser.Marshall(document, nil, 0)
}

/*
 * Simple nested document serialization
 */
func TestSimpleNestedDocumentSerialization(t *testing.T) {
	// Expected buffer from serialization
	var expectedBuffer = []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}
	document := NewDocument()
	subdocument := NewDocument()
	subdocument.Add("int", int32(10))
	document.Add("string", "hello world")
	document.Add("doc", subdocument)

	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	// // Serialize tests
	SerializeTest(t, &T2{"hello world", &T1{10}}, expectedBuffer)
	SerializeTest(t, document, expectedBuffer)

	// De serializing tests
	DeserializeTest(t, expectedBuffer, NewDocument(), document)
	DeserializeTest(t, expectedBuffer, &T2{}, &T2{"hello world", &T1{10}})
}

/*
 * Replacement GetBSON test
 */
type GetBSONT1 struct {
	Int int32 `bson:"int,omitempty"`
}

func (p *GetBSONT1) GetBSON() (interface{}, error) {
	return &GetBSONT2{"hello world"}, nil
}

type GetBSONT2 struct {
	String string `bson:"string,omitempty"`
}

func TestSimpleGetBSONReplacement(t *testing.T) {
	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}

	// Serialize tests
	SerializeTest(t, &GetBSONT1{10}, expectedBuffer)

	// Deserialize test
	DeserializeTest(t, expectedBuffer, &GetBSONT1{}, &GetBSONT1{})
}

type T struct{}
type T1 struct {
	Int int32 `bson:"int,omitempty"`
}
type T2 struct {
	String string `bson:"string,omitempty"`
}

/*
 * Simple Type serializations
 */
func TestTypeSerializations(t *testing.T) {
	type T struct{}
	// Expected buffer from serialization
	var expectedBuffer = []byte{5, 0, 0, 0, 0}
	// Serialize tests
	SerializeTest(t, NewDocument(), expectedBuffer)
	SerializeTest(t, &T{}, expectedBuffer)
}

func TestDocumentWithInt32Serialization(t *testing.T) {
	type T struct {
		Int int32 `bson:"int,omitempty"`
	}
	// Expected buffer from serialization
	var expectedBuffer = []byte{14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	// Create document
	document := NewDocument()
	document.Add("int", int32(10))

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T{10}, expectedBuffer)
}

func TestSimpleStringSerialization(t *testing.T) {
	type T struct {
		String string `bson:"string,omitempty"`
	}
	// Expected buffer from serialization
	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}
	document := NewDocument()
	document.Add("string", "hello world")

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T{"hello world"}, expectedBuffer)
}

func TestSimpleStringAndIntSerialization(t *testing.T) {
	type T struct {
		String string `bson:"string,omitempty"`
		Int    int32  `bson:"int,omitempty"`
	}
	// Expected buffer from serialization
	var expectedBuffer = []byte{38, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	document := NewDocument()
	document.Add("string", "hello world")
	document.Add("int", int32(10))

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T{"hello world", 10}, expectedBuffer)
}

func TestSimpleArraySerialization(t *testing.T) {
	var expectedBuffer = []byte{35, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 23, 0, 0, 0, 2, 48, 0, 2, 0, 0, 0, 97, 0, 2, 49, 0, 2, 0, 0, 0, 98, 0, 0, 0}
	document := NewDocument()
	array := make([]interface{}, 0)
	array = append(array, "a")
	array = append(array, "b")
	document.Add("array", array)

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)

	// t.Logf("[%v]", len(bson))
	// t.Logf("[%v]", bson)
	// t.Logf("[%v]", expectedBuffer)

	// if err != nil {
	// 	t.Fatalf("Failed to create bson document %v", err)
	// }

	// if len(bson) != len(expectedBuffer) {
	// 	t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	// }

	// if bytes.Compare(bson, expectedBuffer) != 0 {
	// 	t.Fatalf("Illegal BSON returned")
	// }

	// // Deserialize the object
	// obj, err := Deserialize(expectedBuffer)
	// if err != nil {
	// 	t.Fatalf("Failed to deserialize the bson array")
	// }

	// validateObjectSize(t, obj, 1)
	// a, _ := obj.FieldAsArray("array")
	// validateString(t, a[0], "a")
	// validateString(t, a[1], "b")
}
