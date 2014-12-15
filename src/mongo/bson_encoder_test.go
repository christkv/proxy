package mongo

import (
	"bytes"
	// "fmt"
	// "strings"
	// "reflect"
	"testing"
	// "time"
)

func SerializeTest(t *testing.T, doc interface{}, expectedBuffer []byte) {
	// bson := ()
	// // Validate the size of the bson array
	// size, _ := bson.CalculateObjectSize(reflect.ValueOf(doc))
	// if size != len(expectedBuffer) {
	// 	t.Errorf("size comparison failed %v != %v", size, len(expectedBuffer))
	// }

	// Serialize the document allowing self allocation of buffer
	b, err := Marshall(doc, nil, 0)

	t.Logf("%s", b)

	// Ensure the buffers match
	if err != nil || len(b) != len(expectedBuffer) || bytes.Compare(b, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned \nexp: %v:%v\ngot: %v:%v", expectedBuffer, len(expectedBuffer), b, len(b))
	}

	// Serialize into pre-allocated buffer
	b = make([]byte, len(expectedBuffer))
	// Serialize the document
	b, err = Marshall(doc, b, 0)
	// Ensure the buffers match
	if err != nil || len(b) != len(expectedBuffer) || bytes.Compare(b, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned \nexp: %v:%v\ngot: %v:%v", expectedBuffer, len(expectedBuffer), b, len(b))
	}
}

// func DeserializeTest(t *testing.T, b []byte, empty interface{}, expected interface{}) {
// 	bson := NewBSON()
// 	// Deserialize the data into the type
// 	err := bson.Deserialize(b, empty)
// 	if err != nil {
// 		t.Errorf("[%v] Failed to deserialize %v into type %v", err, b, expected)
// 	}

// 	// Check if this is a document
// 	switch doc := empty.(type) {
// 	case *Document:
// 		switch doc1 := expected.(type) {
// 		case Document:
// 			if doc1.Equal(doc) == false {
// 				t.Errorf("failed to deserialize document correctly 3")
// 			}
// 		case *Document:
// 			if doc1.Equal(doc) == false {
// 				t.Errorf("failed to deserialize document correctly 4")
// 			}
// 		}
// 	default:
// 		if reflect.DeepEqual(empty, expected) == false {
// 			t.Errorf("failed to deserialize struct correctly 5")
// 		}
// 	}
// }

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
	// SerializeTest(t, &T2{"hello world", &T1{10}}, expectedBuffer)
	SerializeTest(t, document, expectedBuffer)

	// // De serializing tests
	// DeserializeTest(t, expectedBuffer, NewDocument(), document)
	// DeserializeTest(t, expectedBuffer, &T2{}, &T2{"hello world", &T1{10}})
}
