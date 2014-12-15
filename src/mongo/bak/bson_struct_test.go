package bak

// import (
// 	"bytes"
// 	"log"
// 	"testing"
// )

// func TestSimpleStructSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{5, 0, 0, 0, 0}

// 	type Test struct {
// 	}

// 	// Create instance
// 	obj := &Test{}

// 	// Serialize the struct
// 	bson, err := Serialize(obj, nil, 0)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document")
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}
// }

// func TestSimpleStructSerializationWithInt(t *testing.T) {
// 	var expectedBuffer = []byte{14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}

// 	type Test struct {
// 		Int int32 `bson:"int,omitempty"`
// 	}

// 	// Create instance
// 	obj := &Test{10}

// 	// Serialize the struct
// 	bson, err := Serialize(obj, nil, 0)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document")
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}
// }

// func TestSimpleStructStringSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}

// 	type Test struct {
// 		String string `bson:"string,omitempty"`
// 	}

// 	// Create instance
// 	obj := &Test{"hello world"}

// 	// Serialize the struct
// 	bson, err := Serialize(obj, nil, 0)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document")
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}
// }

// func TestSimpleStructStringAndIntSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{38, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}

// 	type Test struct {
// 		String string `bson:"string,omitempty"`
// 		Int    int32  `bson:"int,omitempty"`
// 	}

// 	// Create instance
// 	obj := &Test{"hello world", int32(10)}

// 	// Serialize the struct
// 	bson, err := Serialize(obj, nil, 0)

// 	log.Printf("================================\n%v\n%v", bson, err)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document")
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}
// }

// func TestSimpleStructNestedDocumentSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}

// 	type Test1 struct {
// 		Int int32 `bson:"int,omitempty"`
// 	}

// 	type Test struct {
// 		String string `bson:"string,omitempty"`
// 		Doc    *Test1 `bson:"doc,omitempty"`
// 	}

// 	// Create instance
// 	obj := &Test{"hello world", &Test1{int32(10)}}

// 	// Serialize the struct
// 	bson, err := Serialize(obj, nil, 0)

// 	log.Printf("================================\n%v\n%v", bson, err)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document")
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// document := NewDocument()
// 	// subdocument := NewDocument()
// 	// subdocument.Add("int", int32(10))
// 	// document.Add("string", "hello world")
// 	// document.Add("doc", subdocument)
// 	// bson, err := Serialize(document, nil, 0)

// 	// t.Logf("[%v]", len(bson))
// 	// t.Logf("[%v]", bson)
// 	// t.Logf("[%v]", expectedBuffer)

// 	// if err != nil {
// 	// 	t.Fatalf("Failed to create bson document %v", err)
// 	// }

// 	// if len(bson) != len(expectedBuffer) {
// 	// 	t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	// }

// 	// if bytes.Compare(bson, expectedBuffer) != 0 {
// 	// 	t.Fatalf("Illegal BSON returned")
// 	// }

// 	// // Deserialize the object
// 	// obj, err := Deserialize(expectedBuffer)
// 	// if err != nil {
// 	// 	t.Fatalf("Failed to deserialize the bson array")
// 	// }

// 	// validateObjectSize(t, obj, 2)
// 	// validateStringField(t, obj, "string", "hello world")
// 	// validateIntField(t, subDocument(obj, "doc"), "int", 10)
// }
