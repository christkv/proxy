package bak

// import (
// 	"gopkg.in/mgo.v2/bson"
// 	"reflect"
// 	"testing"
// )

// func BenchmarkNestedDocumentSerializationMGO(b *testing.B) {
// 	document := NewDocument()
// 	subdocument := NewDocument()
// 	subdocument.Add("int", int32(10))
// 	document.Add("string", "hello world")
// 	document.Add("doc", subdocument)

// 	type T1 struct {
// 		Int int32 `bson:"int,omitempty"`
// 	}

// 	type T2 struct {
// 		String string `bson:"string,omitempty"`
// 		Doc    *T1    `bson:"doc,omitempty"`
// 	}

// 	obj := &T2{"hello world", &T1{10}}

// 	for n := 0; n < b.N; n++ {
// 		bson.Marshal(obj)
// 	}
// }

// func BenchmarkNestedDocumentSerialization(b *testing.B) {
// 	document := NewDocument()
// 	subdocument := NewDocument()
// 	subdocument.Add("int", int32(10))
// 	document.Add("string", "hello world")
// 	document.Add("doc", subdocument)

// 	type T1 struct {
// 		Int int32 `bson:"int,omitempty"`
// 	}

// 	type T2 struct {
// 		String string `bson:"string,omitempty"`
// 		Doc    *T1    `bson:"doc,omitempty"`
// 	}

// 	parser := NewBSON()
// 	obj := &T2{"hello world", &T1{10}}
// 	value := reflect.ValueOf(obj)
// 	size, _ := parser.CalculateObjectSize(value)
// 	buffer := make([]byte, size)

// 	for n := 0; n < b.N; n++ {
// 		parser.Serialize(obj, buffer, 0)
// 	}
// }

// func BenchmarkNestedDocumentSerializationNoPreAllocation(b *testing.B) {
// 	document := NewDocument()
// 	subdocument := NewDocument()
// 	subdocument.Add("int", int32(10))
// 	document.Add("string", "hello world")
// 	document.Add("doc", subdocument)

// 	type T1 struct {
// 		Int int32 `bson:"int,omitempty"`
// 	}

// 	type T2 struct {
// 		String string `bson:"string,omitempty"`
// 		Doc    *T1    `bson:"doc,omitempty"`
// 	}

// 	parser := NewBSON()
// 	obj := &T2{"hello world", &T1{10}}

// 	for n := 0; n < b.N; n++ {
// 		parser.Serialize(obj, nil, 0)
// 	}
// }
