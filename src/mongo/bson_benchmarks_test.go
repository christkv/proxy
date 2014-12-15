package mongo

import (
	"testing"
)

func BenchmarkNestedDocumentSerializationStruct(b *testing.B) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	// parser := NewBSON()
	obj := &T2{"hello world", &T1{10}}
	bson := NewBSON()

	for n := 0; n < b.N; n++ {
		bson.Marshall(obj, nil, 0)
	}
}

func BenchmarkNestedDocumentSerializationDocument(b *testing.B) {
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

	bson := NewBSON()
	for n := 0; n < b.N; n++ {
		bson.Marshall(document, nil, 0)
	}
}
