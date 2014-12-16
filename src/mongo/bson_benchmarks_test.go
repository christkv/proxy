package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

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

	parser := NewBSON()
	for n := 0; n < b.N; n++ {
		parser.Marshall(document, nil, 0)
	}
}

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
	parser := NewBSON()

	for n := 0; n < b.N; n++ {
		parser.Marshall(obj, nil, 0)
	}
}

func BenchmarkNestedDocumentSerializationMGO(b *testing.B) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	obj := &T2{"hello world", &T1{10}}

	for n := 0; n < b.N; n++ {
		bson.Marshal(obj)
	}
}

func BenchmarkNestedDocumentSerializationStructOverFlow64bytes(b *testing.B) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	// parser := NewBSON()
	obj := &T2{"hello world hello world hello world hello world hello world hello world", &T1{10}}
	parser := NewBSON()

	for n := 0; n < b.N; n++ {
		parser.Marshall(obj, nil, 0)
	}
}

func BenchmarkNestedDocumentSerializationMGOverflow64Bytes(b *testing.B) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	obj := &T2{"hello world hello world hello world hello world hello world hello world", &T1{10}}

	for n := 0; n < b.N; n++ {
		bson.Marshal(obj)
	}
}

func BenchmarkNestedDocumentDeserialization(b *testing.B) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	parser := NewBSON()
	data := []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}
	obj := &T2{}

	for n := 0; n < b.N; n++ {
		parser.Unmarshal(data, obj)
	}
}

func BenchmarkNestedDocumentDeserializationMGO(b *testing.B) {
	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	data := []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}
	obj := &T2{}

	for n := 0; n < b.N; n++ {
		bson.Unmarshal(data, obj)
	}
}

type GetBSONBenchmarkT1 struct {
	Int int32 `bson:"int,omitempty"`
}

func (p *GetBSONBenchmarkT1) GetBSON() (interface{}, error) {
	return &GetBSONBenchmarkT2{"hello world"}, nil
}

type GetBSONBenchmarkT2 struct {
	String string `bson:"string,omitempty"`
}

func BenchmarkGetBSONSerialization(b *testing.B) {
	doc := &GetBSONT1{10}
	parser := NewBSON()

	for n := 0; n < b.N; n++ {
		parser.Marshall(doc, nil, 0)
	}
}

func BenchmarkGetBSONSerializationMGO(b *testing.B) {
	doc := &GetBSONT1{10}

	for n := 0; n < b.N; n++ {
		bson.Marshal(doc)
	}
}
