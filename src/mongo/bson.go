package mongo

import (
	"reflect"
	"strings"
)

type bsonType byte

const (
	bsonString   bsonType = 0x02
	bsonDocument bsonType = 0x03
	bsonInt32    bsonType = 0x10
)

type ObjectId struct {
	Id []byte
}

type Binary struct {
	SubType byte
	Data    []byte
}

type Javascript struct {
	Code string
}

type JavascriptWScope struct {
	Code  string
	Scope *Document
}

type Date struct {
	Value int64
}

type RegExp struct {
	Pattern string
	Options string
}

type Timestamp struct {
	Value int64
}

type Min struct {
}

type Max struct {
}

type DBPointer struct {
}

type TypeInfos struct {
	Types map[string]*TypeInfo
}

type TypeInfo struct {
	Fields        map[string]*FieldInfo
	FieldsByIndex []*FieldInfo
	NumberOfField int
}

type FieldInfo struct {
	Name         string
	MetaDataName string
}

func writeU32(buffer []byte, index int, value uint32) {
	buffer[index+3] = byte((value >> 24) & 0xff)
	buffer[index+2] = byte((value >> 16) & 0xff)
	buffer[index+1] = byte((value >> 8) & 0xff)
	buffer[index] = byte(value & 0xff)
}

func writeU64(buffer []byte, index int, value uint64) {
	buffer[index+7] = byte((value >> 56) & 0xff)
	buffer[index+6] = byte((value >> 48) & 0xff)
	buffer[index+5] = byte((value >> 40) & 0xff)
	buffer[index+4] = byte((value >> 32) & 0xff)
	buffer[index+3] = byte((value >> 24) & 0xff)
	buffer[index+2] = byte((value >> 16) & 0xff)
	buffer[index+1] = byte((value >> 8) & 0xff)
	buffer[index] = byte(value & 0xff)
}

func readUInt64(buffer []byte, index int) uint64 {
	return (uint64(buffer[index]) << 0) |
		(uint64(buffer[index+1]) << 8) |
		(uint64(buffer[index+2]) << 16) |
		(uint64(buffer[index+3]) << 24) |
		(uint64(buffer[index+4]) << 32) |
		(uint64(buffer[index+5]) << 40) |
		(uint64(buffer[index+6]) << 48) |
		(uint64(buffer[index+7]) << 56)
}

func readUInt32(buffer []byte, index int) uint32 {
	return (uint32(buffer[index]) << 0) |
		(uint32(buffer[index+1]) << 8) |
		(uint32(buffer[index+2]) << 16) |
		(uint32(buffer[index+3]) << 24)
}

type Getter interface {
	GetBSON() (interface{}, error)
}

type BSON struct {
	typeInfos *TypeInfos
}

func NewBSON() *BSON {
	return &BSON{&TypeInfos{make(map[string]*TypeInfo)}}
}

func parseTypeInformation(typeInfos *TypeInfos, value reflect.Value) *TypeInfo {
	// We have a pointer get the underlying value
	if value.Type().Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Reuse type information if already present
	if typeInfos.Types[value.Type().Name()] != nil {
		return typeInfos.Types[value.Type().Name()]
	}

	// Get the number of fields
	numberOfFields := value.NumField()

	// Create typeInfo box
	typeInfo := TypeInfo{}
	// Pre-allocate a map with the entries we need
	typeInfo.Fields = make(map[string]*FieldInfo, numberOfFields*2)
	typeInfo.FieldsByIndex = make([]*FieldInfo, numberOfFields)
	typeInfo.NumberOfField = numberOfFields

	// Iterate over all the fields and collect the metadata
	for index := 0; index < numberOfFields; index++ {
		// Get the field information
		fieldType := value.Type().Field(index)
		// Get the field name
		key := fieldType.Name
		// Get the tag for the field
		tag := fieldType.Tag.Get("bson")

		// Split the tag into parts
		parts := strings.Split(tag, ",")

		// Override the key if the metadata has one
		if len(parts) > 0 && parts[0] != "" {
			key = parts[0]
		}

		// Create a new fieldInfo instance
		fieldInfo := FieldInfo{fieldType.Name, key}
		// Add to the map
		typeInfo.Fields[fieldType.Name] = &fieldInfo
		typeInfo.Fields[key] = &fieldInfo
		typeInfo.FieldsByIndex[index] = &fieldInfo
	}

	// Save type
	typeInfos.Types[value.Type().Name()] = &typeInfo
	// Return the type information
	return &typeInfo
}
