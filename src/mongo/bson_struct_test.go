package mongo

import (
	"bytes"
	"log"
	"testing"
)

func TestSimpleStructSerialization(t *testing.T) {
	var expectedBuffer = []byte{5, 0, 0, 0, 0}

	type Test struct {
	}

	// Create instance
	obj := &Test{}

	// Serialize the struct
	bson, err := Serialize(obj, nil, 0)

	if err != nil {
		t.Fatalf("Failed to create bson document")
	}

	if len(bson) != len(expectedBuffer) {
		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned")
	}
}

func TestSimpleStructSerializationWithInt(t *testing.T) {
	var expectedBuffer = []byte{14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}

	type Test struct {
		Int int32 `bson:"int,omitempty"`
	}

	// Create instance
	obj := &Test{10}

	// Serialize the struct
	bson, err := Serialize(obj, nil, 0)

	log.Printf("================================\n%v\n%v", bson, err)

	if err != nil {
		t.Fatalf("Failed to create bson document")
	}

	if len(bson) != len(expectedBuffer) {
		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned")
	}
}

func TestSimpleStructStringSerialization(t *testing.T) {
	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}

	type Test struct {
		String string `bson:"string,omitempty"`
	}

	// Create instance
	obj := &Test{"hello world"}

	// Serialize the struct
	bson, err := Serialize(obj, nil, 0)

	log.Printf("================================\n%v\n%v", bson, err)

	if err != nil {
		t.Fatalf("Failed to create bson document")
	}

	if len(bson) != len(expectedBuffer) {
		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned")
	}
}
