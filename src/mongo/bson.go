package mongo

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	// "math"
	"reflect"
	"strconv"
	"strings"
	// "time"
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

func Parse(bson []byte) (interface{}, error) {
	return nil, nil
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

func calculateElementSize(value reflect.Value) (int, error) {
	size := 0
	// log.Printf("=========== calculateElementSize :: %s", value.Kind())

	// Store original type
	// originalValue := value
	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Switch on the value
	switch value.Kind() {
	case reflect.Int32, reflect.Uint32:
		// log.Printf("got Int32")
		size = size + 4
	case reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64:
		// log.Printf("got Int64")
		size = size + 8
	case reflect.String:
		// log.Printf("got String")
		size = 4 + len(value.String()) + 1
	case reflect.Struct:
		// log.Printf("got Struct")
		// log.Printf("############################# STRUCT HUH")
		// Switch on the type
		switch value.Interface().(type) {
		case ObjectId:
			size = size + 12
		case Document:
			// log.Printf("got Document")
			elementSize, err := CalculateObjectSize(value)
			if err != nil {
				return size, err
			}

			size = size + elementSize
		default:
			// log.Printf("got Document")

			elementSize, err := CalculateObjectSize(value)
			if err != nil {
				return size, err
			}

			size = size + elementSize
		}
	default:
		return 0, errors.New(fmt.Sprintf("BSON struct type %T not supported for serialization", value))
	}

	return size, nil

	// // Serialize the document
	// switch element := elem.(type) {
	// default:
	// 	return size, errors.New(fmt.Sprintf("unsupported type %T", element))
	// case reflect.Value:
	// 	// // Get type of
	// 	// typeof := reflect.ValueOf(elem)
	// 	// // If we have a pointer get actual element
	// 	// if elem.Kind() == reflect.Ptr {
	// 	// 	typeof = typeof.Elem()
	// 	// }

	// 	switch element.Kind() {
	// 	case reflect.Int32, reflect.Uint32:
	// 		size = size + 4
	// 	case reflect.String:
	// 		size = 4 + len(element.String()) + 1
	// 	case reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64:
	// 		size = size + 8
	// 	case reflect.Struct:
	// 		elementSize, err := CalculateObjectSize(element)
	// 		if err != nil {
	// 			return size, err
	// 		}

	// 		size = size + elementSize
	// 	default:
	// 		// switch reflect.ValueOf(elem).Type().Kind() {
	// 		// case reflect.Struct:
	// 		// 	log.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^ STRUCT")
	// 		// 	elementSize, err := CalculateObjectSize(reflect.ValueOf(elem))
	// 		// 	if err != nil {
	// 		// 		return size, err
	// 		// 	}

	// 		// 	size = size + elementSize
	// 		// 	return size, nil
	// 		// }

	// 		// typeof := reflect.ValueOf(elem).Type()
	// 		// // if typeof.Kind() == reflect.Ptr {
	// 		// log.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^ POINTER %v", typeof)
	// 		// 	// 	typeof = typeof.Elem()
	// 		// }

	// 		return size, errors.New(fmt.Sprintf("unsupported type %T", elem))
	// 	}
	// case []interface{}:
	// 	elementSize, err := calculateArraySize(element)
	// 	if err != nil {
	// 		return size, err
	// 	}

	// 	size = size + elementSize
	// case *Document:
	// 	elementSize, err := CalculateObjectSize(element)
	// 	if err != nil {
	// 		return size, err
	// 	}

	// 	size = size + elementSize
	// case string:
	// 	size = size + 4 + len(element) + 1
	// case int32, uint32:
	// 	size = size + 4
	// case int64, uint64, float32, float64, *Date, *Timestamp, time.Time, *time.Time:
	// 	size = size + 8
	// case nil:
	// 	size = size
	// case *Binary:
	// 	size = size + 4 + 1 + len(element.Data)
	// case bool:
	// 	size = size + 1
	// case []byte:
	// 	size = size + 4 + 1 + len(element)
	// case *ObjectId:
	// 	size = size + 12
	// case *RegExp:
	// 	size = size + len(element.Pattern) + 1 + len(element.Options) + 1
	// case *Javascript:
	// 	size = size + len(element.Code) + 4 + 1
	// case *JavascriptWScope:
	// 	elementSize, err := CalculateObjectSize(element.Scope)
	// 	if err != nil {
	// 		return size, err
	// 	}

	// 	size = size + len(element.Code) + elementSize + 4 + 1 + 4
	// case Min, *Min, Max, *Max:
	// 	size = size
	// }

	// return size, nil
}

func calculateArraySize(array []interface{}) (int, error) {
	size := 5

	// Iterate over all the key values
	for index, value := range array {
		indexStr := strconv.Itoa(index)
		// Add the key size
		size = size + len(indexStr) + 1 + 1
		// Add the size of the actual element
		elementSize, err := calculateElementSize(reflect.ValueOf(value))

		if err != nil {
			return size, err
		}

		size = size + elementSize
	}

	return size, nil
}

func packString(buffer []byte, originalIndex int, index int, value reflect.Value) int {
	// Get the string
	str := value.String()
	// Set the type
	buffer[originalIndex] = 0x02
	// Get the string bytes
	stringBytes := []byte(str)
	// Write the string length
	writeU32(buffer, index+1, uint32(len(stringBytes)+1))
	// Copy the string bytes to the bson buffer
	copy(buffer[index+5:], stringBytes[:])
	// Set end 0 byte
	buffer[index+5+len(stringBytes)] = 0x00
	// Return new index position
	return index + 4 + len(stringBytes) + 1 + 1
}

func packInt32(buffer []byte, originalIndex int, index int, value uint32) int {
	buffer[originalIndex] = 0x10
	writeU32(buffer, index+1, value)
	return index + 5
}

func packElement(key string, value reflect.Value, buffer []byte, index int) (int, error) {
	log.Printf("packElement %v with value %v of kind %v", key, value, value.Kind())
	strbytes := []byte(key)
	// Save a pointer to the first byte index
	originalIndex := index
	// Skip type
	index = index + 1
	// Copy the string into the buffer
	copy(buffer[index:], strbytes[0:])
	// Null terminate the string
	buffer[index+len(strbytes)] = 0

	// Update the index with the field length
	index = index + len(strbytes)

	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Reflect on the type
	switch value.Kind() {
	case reflect.String:
		log.Printf("packElement ================ got string")
		index = packString(buffer, originalIndex, index, value)
	case reflect.Int32:
		log.Printf("packElement ================ got int32 %v", value.Int())
		buffer[originalIndex] = 0x10
		writeU32(buffer, index+1, uint32(value.Int()))
		index = index + 5
	case reflect.Struct:
		switch value.Interface().(type) {
		case Document:
			log.Printf("packElement ================ got document")
			// Set the type of be document
			buffer[originalIndex] = 0x03
			// Get the final values
			in, err := serializeObject(buffer, index+1, value)
			// Serialize the object
			return in + 1, err

		default:
			log.Printf("packElement ================ got struct")
			// Set the type of be document
			buffer[originalIndex] = 0x03
			// Get the final values
			in, err := serializeObject(buffer, index+1, value)
			// Serialize the object
			return in + 1, err
		}
	default:
		return index, errors.New(fmt.Sprintf("could not recognize the type %v", value.Kind()))
	}

	// Return the index
	return index, nil
	// strbytes := []byte(key)
	// // Save a pointer to the first byte index
	// originalIndex := index
	// // Skip type
	// index = index + 1
	// // Copy the string into the buffer
	// copy(buffer[index:], strbytes[0:])
	// // Null terminate the string
	// buffer[index+len(strbytes)] = 0

	// // Update the index with the field length
	// index = index + len(strbytes)

	// log.Printf("############################### PACK")
	// // Determine the type
	// switch element := value.(type) {
	// default:
	// 	return index, errors.New(fmt.Sprintf("unsupported type %T", element))
	// case reflect.Value:
	// 	log.Printf("REFLECTED VALUE")
	// 	switch element.Kind() {
	// 	case reflect.Int32:
	// 		log.Printf("reflected int32 serialize %v", uint32(element.Int()))
	// 		index = packInt32(buffer, originalIndex, index, uint32(element.Int()))
	// 	case reflect.String:
	// 		log.Printf("reflected string serialize %v", element.String())
	// 		index = packString(buffer, originalIndex, index, element.String())
	// 	}
	// case int32:
	// 	log.Printf("int32 serialize")
	// 	buffer[originalIndex] = 0x10
	// 	writeU32(buffer, index+1, uint32(element))
	// 	index = index + 5
	// case uint32:
	// 	log.Printf("uint32 serialize")
	// 	buffer[originalIndex] = 0x10
	// 	writeU32(buffer, index+1, element)
	// 	index = index + 5
	// case int64:
	// 	log.Printf("int64 serialize")
	// 	buffer[originalIndex] = 0x12
	// 	writeU64(buffer, index+1, uint64(element))
	// 	index = index + 9
	// case uint64:
	// 	log.Printf("uint64 serialize")
	// 	buffer[originalIndex] = 0x12
	// 	writeU64(buffer, index+1, element)
	// 	index = index + 9
	// case bool:
	// 	log.Printf("bool serialize")
	// 	buffer[originalIndex] = 0x08
	// 	if element {
	// 		buffer[index+1] = 0x01
	// 	} else {
	// 		buffer[index+1] = 0x00
	// 	}
	// 	index = index + 2
	// case float32:
	// 	log.Printf("float32 serialize")
	// 	// int64(math.Float64bits(v)
	// 	buffer[originalIndex] = 0x01
	// 	// Get reflection of the value
	// 	reflectType := reflect.ValueOf(element)
	// 	// Convert 32 bit float to string
	// 	floatString := strconv.FormatFloat(reflectType.Float(), 'g', -1, 32)
	// 	// Parse string as 64bit float
	// 	value, _ := strconv.ParseFloat(floatString, 64)
	// 	// Write the float as an uint64
	// 	writeU64(buffer, index+1, math.Float64bits(value))
	// 	index = index + 9
	// case float64:
	// 	log.Printf("float64 serialize")
	// 	buffer[originalIndex] = 0x01
	// 	writeU64(buffer, index+1, math.Float64bits(element))
	// 	index = index + 9
	// case *RegExp:
	// 	log.Printf("regexp serialize")
	// 	buffer[originalIndex] = 0x0b
	// 	copy(buffer[index+1:], []byte(element.Pattern))
	// 	index = index + 1 + len(element.Pattern)
	// 	buffer[index] = 0x00
	// 	copy(buffer[index+1:], []byte(element.Options))
	// 	index = index + 1 + len(element.Options) + 1
	// case *Timestamp:
	// 	log.Printf("timestamp serialize")
	// 	buffer[originalIndex] = 0x11
	// 	writeU64(buffer, index+1, uint64(element.Value))
	// 	index = index + 9
	// case *Date:
	// 	log.Printf("date serialize")
	// 	buffer[originalIndex] = 0x09
	// 	writeU64(buffer, index+1, uint64(element.Value))
	// 	index = index + 9
	// case *time.Time:
	// 	log.Printf("*time.Time serialize")
	// 	buffer[originalIndex] = 0x09
	// 	writeU64(buffer, index+1, uint64(element.Unix()))
	// 	index = index + 9
	// case time.Time:
	// 	log.Printf("time.Time serialize")
	// 	buffer[originalIndex] = 0x09
	// 	writeU64(buffer, index+1, uint64(element.Unix()))
	// 	index = index + 9
	// case nil:
	// 	log.Printf("nil serialize")
	// 	buffer[originalIndex] = 0x0a
	// 	index = index + 1
	// case string:
	// 	log.Printf("string serialize")
	// 	index = packString(buffer, originalIndex, index, element)
	// case []interface{}:
	// 	log.Printf("array serialize %v", index)
	// 	// Set the type of be document
	// 	buffer[originalIndex] = 0x04
	// 	// Get the final values
	// 	in, err := serializeArray(buffer, index+1, element)
	// 	// Serialize the object
	// 	return in + 1, err
	// case *Document:
	// 	log.Printf("document serialize %v", index)
	// 	// Set the type of be document
	// 	buffer[originalIndex] = 0x03
	// 	// Get the final values
	// 	in, err := serializeObject(buffer, index+1, element)
	// 	// Serialize the object
	// 	return in + 1, err
	// case *ObjectId:
	// 	log.Printf("objectid serialize %v", index)
	// 	if len(element.Id) != 12 {
	// 		return 0, errors.New("ObjectId must be a 12 byte array")
	// 	}

	// 	// Set the type of be document
	// 	buffer[originalIndex] = 0x07
	// 	copy(buffer[index+1:], element.Id)
	// 	return index + len(element.Id) + 1, nil
	// case []byte:
	// 	log.Printf("[]byte serialize %v", index)
	// 	// Set the type of be document
	// 	buffer[originalIndex] = 0x05
	// 	// Set the size of the binary
	// 	writeU32(buffer, index+1, uint32(len(element)))
	// 	buffer[index+5] = 0x00
	// 	// Write binary
	// 	copy(buffer[index+6:], element[:])
	// 	// Return the length
	// 	return index + len(element) + 5 + 1, nil
	// case *Binary:
	// 	log.Printf("binary serialize %v", index)
	// 	// Set the type of be document
	// 	buffer[originalIndex] = 0x05
	// 	// Set the size of the binary
	// 	writeU32(buffer, index+1, uint32(len(element.Data)))
	// 	buffer[index+5] = element.SubType
	// 	// Write binary
	// 	copy(buffer[index+6:], element.Data[:])
	// 	// Return the length
	// 	return index + len(element.Data) + 5 + 1, nil
	// case *Javascript:
	// 	log.Printf("javascript no scope serialize %v", index)
	// 	// Set the type of be document
	// 	buffer[originalIndex] = 0x0d
	// 	stringBytes := []byte(element.Code)
	// 	writeU32(buffer, index+1, uint32(len(stringBytes)+1))
	// 	copy(buffer[index+5:], stringBytes[:])
	// 	buffer[index+5+len(stringBytes)] = 0x00
	// 	index = index + 4 + len(stringBytes) + 1 + 1
	// case *JavascriptWScope:
	// 	log.Printf("javascript scope serialize %v", index)
	// 	buffer[originalIndex] = 0x0f
	// 	stringBytes := []byte(element.Code)
	// 	// Skip the length
	// 	lengthIndex := index + 1
	// 	index = index + 4

	// 	// Write javascript string
	// 	writeU32(buffer, index+1, uint32(len(stringBytes)+1))
	// 	copy(buffer[index+5:], stringBytes[:])
	// 	buffer[index+5+len(stringBytes)] = 0x00
	// 	index = index + 4 + len(stringBytes) + 1

	// 	// Serialize the scope
	// 	in, err := serializeObject(buffer, index+1, element.Scope)
	// 	index = index + in

	// 	// Write the length
	// 	writeU32(buffer, lengthIndex, uint32(in-lengthIndex+1))
	// 	return in + 1, err
	// case *Min:
	// 	log.Printf("min serialize %v", index)
	// 	buffer[originalIndex] = 0xFF
	// 	index = index + 1
	// case *Max:
	// 	log.Printf("max serialize %v", index)
	// 	buffer[originalIndex] = 0x7F
	// 	index = index + 1
	// }

	// // Return the index
	// return index, nil
}

func serializeObject(buffer []byte, index int, value reflect.Value) (int, error) {
	i := index + 4

	log.Printf("Serialize ========================== 0")

	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Switch on the value
	switch value.Kind() {
	case reflect.Struct:
		log.Printf("Serialize ========================== 1")
		switch value.Interface().(type) {
		case Document:
			// Cast to document type
			switch doc := value.Interface().(type) {
			case Document:
				log.Printf("GOT *Document, %v", doc.fields)

				// Iterate over all the key values
				for _, key := range doc.fields {
					log.Printf("key :: %v", key)
					// Get the value
					fieldValue := doc.document[key]
					// Add the size of the actual element
					in, err := packElement(key, reflect.ValueOf(fieldValue), buffer, i)
					if err != nil {
						return i, err
					}

					i = in
				}
			default:
				return i, errors.New(fmt.Sprintf("BSON struct type %T not supported for serialization", doc))
			}
		default:
			log.Printf("Serialize ========================== 2")
			log.Printf("GOT *Struct")
			numberOfField := value.NumField()

			// Let's iterate over all the fields
			for j := 0; j < numberOfField; j++ {
				// Get the field value
				fieldValue := value.Field(j)
				fieldType := value.Type().Field(j)
				// Get the field name
				key := fieldType.Name
				// Get the tag
				tag := fieldType.Tag.Get("bson")
				// Split the tag into parts
				parts := strings.Split(tag, ",")

				// Override the key if the metadata has one
				if len(parts) > 0 && parts[0] != "" {
					key = parts[0]
				}

				// Add the size of the actual element
				in, err := packElement(key, fieldValue, buffer, i)
				if err != nil {
					return i, err
				}

				i = in
			}
		}
	default:
		log.Printf("Serialize ========================== 3")
		log.Printf("BSON struct type %T not supported for serialization", value)
		return i, errors.New(fmt.Sprintf("BSON struct type %T not supported for serialization", value))
	}

	// Final object size
	objectSize := (i - index + 5 - 4)

	// The final object size (the encoded length + terminating 0)
	writeU32(buffer, index, uint32(objectSize))
	// Return no error
	return i, nil

	// i := index + 4

	// switch document := object.(type) {
	// default:
	// 	return 0, errors.New("Unsupported Serialization Mode")
	// case *Document:
	// 	// Iterate over all the key values
	// 	for _, key := range document.fields {
	// 		value := document.document[key]
	// 		in, err := packElement(key, value, buffer, i)
	// 		if err != nil {
	// 			return i, err
	// 		}
	// 		i = in
	// 	}
	// case reflect.Value:
	// 	numberOfField := document.NumField()

	// 	log.Printf("number of fields off struct %v", numberOfField)

	// 	// Let's iterate over all the fields
	// 	for j := 0; j < numberOfField; j++ {
	// 		// Get the field value
	// 		fieldValue := document.Field(j)
	// 		fieldType := document.Type().Field(j)
	// 		// Get the field name
	// 		key := fieldType.Name
	// 		// Get the tag
	// 		tag := fieldType.Tag.Get("bson")
	// 		// Split the tag into parts
	// 		parts := strings.Split(tag, ",")

	// 		// Override the key if the metadata has one
	// 		if len(parts) > 0 && parts[0] != "" {
	// 			key = parts[0]
	// 		}

	// 		log.Printf("serialize field %v of type %v with tag %v at index %v", key, fieldValue, tag, index)

	// 		// Add the size of the actual element
	// 		in, err := packElement(key, fieldValue, buffer, i)
	// 		if err != nil {
	// 			return i, err
	// 		}

	// 		i = in
	// 	}
	// }

	// // Final object size
	// objectSize := (i - index + 5 - 4)

	// // The final object size (the encoded length + terminating 0)
	// writeU32(buffer, index, uint32(objectSize))
	// // Return no error
	// return i, nil
}

// func serializeArray(buffer []byte, index int, array []interface{}) (int, error) {
// 	i := index + 4

// 	// Iterate over all the key values
// 	for index, value := range array {
// 		indexStr := strconv.Itoa(index)
// 		in, err := packElement(indexStr, value, buffer, i)
// 		if err != nil {
// 			return i, err
// 		}
// 		i = in
// 	}

// 	// Final object size
// 	objectSize := (i - index + 5 - 4)

// 	// The final object size (the encoded length + terminating 0)
// 	writeU32(buffer, index, uint32(objectSize))
// 	// Return no error
// 	return i, nil
// }

func CalculateObjectSize(value reflect.Value) (int, error) {
	// Minimum object size
	size := 5

	// originalValue := value
	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// var obj interface{}

	// log.Printf("CalculateObjectSize ========================== 0")

	// Switch on the value
	switch value.Kind() {
	case reflect.Struct:
		switch value.Interface().(type) {
		case ObjectId:
			size = size + 12
		case Document:
			// Cast to document type
			switch doc := value.Interface().(type) {
			case Document:
				// log.Printf("GOT *Document, %v", doc.fields)

				// Iterate over all the key values
				for _, key := range doc.fields {
					// Get the value
					value := doc.document[key]
					// Add the key size
					size = size + len(key) + 1 + 1
					// Add the size of the actual element
					elementSize, err := calculateElementSize(reflect.ValueOf(value))

					if err != nil {
						return size, err
					}

					size = size + elementSize
				}
			default:
				return 0, errors.New(fmt.Sprintf("BSON struct type %T not supported for serialization", doc))
			}
		default:
			numberOfField := value.NumField()

			// Let's iterate over all the fields
			for i := 0; i < numberOfField; i++ {
				// Get the field value
				fieldValue := value.Field(i)
				fieldType := value.Type().Field(i)
				// Get the field name
				key := fieldType.Name
				// Get the tag
				tag := fieldType.Tag.Get("bson")
				// Split the tag into parts
				parts := strings.Split(tag, ",")

				// Override the key if the metadata has one
				if len(parts) > 0 && parts[0] != "" {
					key = parts[0]
				}

				// Add the length of the name of the field
				size = size + len(key) + 1 + 1

				// Add the size of the actual element
				elementSize, err := calculateElementSize(fieldValue)

				if err != nil {
					return size, err
				}

				size = size + elementSize
			}
		}
	default:
		return 0, errors.New(fmt.Sprintf("BSON struct type %T not supported for serialization", value))
	}

	return size, nil
}

func Serialize(obj reflect.Value, bson []byte, offset int) ([]byte, error) {
	// // Minimum object size
	// size := 5

	// // originalValue := value
	// // We have a pointer get the underlying value
	// if value.Kind() == reflect.Ptr {
	// 	value = value.Elem()
	// }

	if bson == nil {
		// Calculate the size of the document
		size, err := CalculateObjectSize(obj)
		if err != nil {
			return nil, err
		}

		bson = make([]byte, size)
		offset = 0
	}

	_, err := serializeObject(bson, offset, obj)
	if err != nil {
		return nil, err
	}

	return bson, nil

	// switch document := obj.(type) {
	// case *Document:
	// 	// We are not using our own buffer to serialize into
	// 	if bson == nil {
	// 		// Calculate the size of the document
	// 		size, err := CalculateObjectSize(reflect.ValueOf(document))
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		log.Printf("size of bson element %v", size)

	// 		// Allocate space
	// 		bson = make([]byte, size)
	// 	}

	// 	// Serialize the object
	// 	_, err := serializeObject(bson[offset:], 0, document)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Return the bson
	// 	return bson, nil
	// default:
	// 	// Get type of
	// 	typeof := reflect.ValueOf(obj)
	// 	// If we have a pointer get actual element
	// 	if typeof.Kind() == reflect.Ptr {
	// 		typeof = typeof.Elem()
	// 	}

	// 	// Check if we have a struct
	// 	switch typeof.Kind() {
	// 	case reflect.Struct:
	// 		// We are not using our own buffer to serialize into
	// 		if bson == nil {
	// 			log.Printf("##########################################################")
	// 			// Calculate the size of the document
	// 			size, err := CalculateObjectSize(reflect.ValueOf(document))
	// 			if err != nil {
	// 				log.Printf("%v", err)
	// 				return nil, err
	// 			}

	// 			log.Printf("size of bson element %v", size)

	// 			// Allocate space
	// 			bson = make([]byte, size)
	// 		}

	// 		// Serialize the object
	// 		_, err := serializeObject(bson[offset:], 0, typeof)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		// Return the bson
	// 		return bson, nil
	// 	}
	// }

	// return nil, errors.New("Unsupported Serialization Mode")
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

// func deserializeArray(bson []byte, index int) ([]interface{}, error) {
// 	// Create document node
// 	array := make([]interface{}, 0)

// 	// Decode the length of the buffer
// 	documentSize := readUInt32(bson, index)
// 	// initialIndex
// 	endIndex := index + int(documentSize)

// 	// Special case of an empty document
// 	if documentSize == 5 {
// 		return array, nil
// 	}

// 	// Skip the size document
// 	index = index + 4

// 	// Start decoding the fields
// 	for index < endIndex-1 {
// 		// Get the bson type
// 		bsonType := bson[index]

// 		// Skip bson type
// 		index = index + 1

// 		// Read the cstring
// 		strindex := bytes.IndexByte(bson[index:], 0x00)

// 		// No 0 byte found error out
// 		if strindex == -1 {
// 			return nil, errors.New("could not decode field name, possibly corrupt bson")
// 		}

// 		// Adjust index with the string length
// 		index = index + strindex + 1

// 		// Switch on type to decode
// 		switch bsonType {
// 		default:
// 			return nil, errors.New(fmt.Sprintf("type [%v] is not a legal bson type", bson[index]))
// 		case 0x00:
// 			index = index + 1
// 			break
// 		case 0x10:
// 			// Read the int
// 			array = append(array, int32(readUInt32(bson, index)))
// 			// Skip the int32 field
// 			index = index + 4
// 		case 0x09:
// 			array = append(array, time.Unix(int64(readUInt64(bson, index)), 0))
// 			// Skip the int64 field
// 			index = index + 8
// 		case 0x11:
// 			array = append(array, &Timestamp{int64(readUInt64(bson, index))})
// 			// Skip the int64 field
// 			index = index + 8
// 		case 0x12:
// 			array = append(array, int64(readUInt64(bson, index)))
// 			// Skip the int64 field
// 			index = index + 8
// 		case 0x01:
// 			array = append(array, math.Float64frombits(readUInt64(bson, index)))
// 			// Skip the int64 field
// 			index = index + 8
// 		case 0x03:
// 			// Read the document size
// 			docSize := int(readUInt32(bson, index))
// 			// Deserialize documents
// 			obj, err := deserializeObject(bson[index:index+docSize], 0)
// 			if err != nil {
// 				return nil, err
// 			}

// 			array = append(array, obj)
// 			// Skip last null byte and size of string
// 			index = index + docSize
// 		case 0x04:
// 			// Read the document size
// 			stringSize := int(readUInt32(bson, index))
// 			// Deserialize documents
// 			obj, err := deserializeArray(bson[index:index+stringSize], 0)
// 			if err != nil {
// 				return nil, err
// 			}
// 			// Add document to array
// 			array = append(array, obj)
// 			// Skip last null byte and size of string
// 			index = index + stringSize
// 		case 0x02:
// 			// Read the string size
// 			stringSize := int(readUInt32(bson, index))
// 			// Skip string size
// 			index = index + 4
// 			// Add document to array
// 			array = append(array, string(bson[index:index+stringSize-1]))
// 			// Skip last null byte and size of string
// 			index = index + stringSize
// 		case 0x05:
// 			// Read the string size
// 			binarySize := int(readUInt32(bson, index))
// 			// Skip string size
// 			index = index + 4
// 			// Get the subtype
// 			subStype := bson[index]
// 			// Skip subtype
// 			index = index + 1
// 			// Add the field
// 			array = append(array, &Binary{subStype, bson[index : index+binarySize]})
// 			// Skip last null byte and size of string
// 			index = index + binarySize
// 		case 0x07:
// 			// Add the field
// 			array = append(array, &ObjectId{bson[index : index+12]})
// 			// Skip last null byte and size of string
// 			index = index + 12
// 		case 0x06:
// 			array = append(array, nil)
// 		case 0x0a:
// 			array = append(array, nil)
// 		case 0x0d:
// 			// Read the string size
// 			stringSize := int(readUInt32(bson, index))
// 			// Skip string size
// 			index = index + 4
// 			// Add the field
// 			array = append(array, &Javascript{string(bson[index : index+stringSize-1])})
// 			// Skip last null byte and size of string
// 			index = index + stringSize
// 		case 0x0F:
// 			// Skip length don't need to decode this
// 			index = index + 4
// 			// Read the string size
// 			stringSize := int(readUInt32(bson, index))
// 			// Skip string size
// 			index = index + 4
// 			// Get the js code string
// 			jsCode := string(bson[index : index+stringSize-1])
// 			// Skip last null byte and size of string
// 			index = index + stringSize
// 			// Read the document size
// 			docSize := int(readUInt32(bson, index))
// 			// Deserialize documents
// 			obj, err := deserializeObject(bson[index:index+docSize], 0)
// 			if err != nil {
// 				return nil, err
// 			}
// 			// Create js object
// 			js := &JavascriptWScope{jsCode, obj}
// 			// Add the field
// 			array = append(array, js)
// 			// Adjust index
// 			index = index + docSize
// 		case 0xff:
// 			array = append(array, &Min{})
// 		case 0x7f:
// 			array = append(array, &Max{})
// 		}
// 	}

// 	// Adjust for last byte
// 	index = index + 1
// 	// Return the document
// 	return array, nil
// }

type TypeInfos struct {
}

type FieldInfo struct {
	Name         string
	MetaDataName string
}

type TypeInfo struct {
	Fields map[string]FieldInfo
}

func parseTypeInformation(value reflect.Value) TypeInfo {
	// Get the number of fields
	numberOfFields := value.NumField()

	// Create typeInfo box
	typeInfo := TypeInfo{}
	// Pre-allocate a map with the entries we need
	typeInfo.Fields = make(map[string]FieldInfo, numberOfFields*2)

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
		typeInfo.Fields[fieldType.Name] = fieldInfo
		typeInfo.Fields[key] = fieldInfo
	}

	return typeInfo
}

func addValueToFieldStruct(fieldName string, obj reflect.Value, value interface{}, isDocument bool) error {
	if isDocument {
		switch t := obj.Interface().(type) {
		case Document:
			t.Add(fieldName, value)
		case *Document:
			t.Add(fieldName, value)
		}
	} else {
		log.Printf("************** addStruct")
		// Get the type info
		typeInfo := parseTypeInformation(obj)
		structFieldName := typeInfo.Fields[fieldName].Name

		// Set the field value on the struct (just set it hard)
		field := obj.FieldByName(structFieldName)

		// We did not find the field on the struct
		if field.Kind() == reflect.Invalid {
			return errors.New(fmt.Sprintf("field %v not found on struct %v", fieldName, obj))
		}

		// Set the field
		field.Set(reflect.ValueOf(value))
	}

	return nil
}

func addDocumentToFieldStruct(fieldName string, obj reflect.Value, isDocument bool) (reflect.Value, error) {
	if isDocument {
		switch t := obj.Interface().(type) {
		case Document:
			// Create a new document
			doc := NewDocument()
			// Add the field
			t.Add(fieldName, doc)
			// Return the new value
			return reflect.ValueOf(doc), nil
		}
	} else {
		// Get the type info
		typeInfo := parseTypeInformation(obj)
		structFieldName := typeInfo.Fields[fieldName].Name

		// Set the field value on the struct (just set it hard)
		// field := obj.FieldByName(structFieldName)

		// Get the type information for the field
		fieldType := obj.Type().FieldByName(structFieldName)
		reflect.New(fieldType.Type)
		// fieldType.Type

		// // We did not find the field on the struct
		// if field.Kind() == reflect.Invalid {
		// 	return errors.New(fmt.Sprintf("field %v not found on struct %v", fieldName, obj))
		// }

	}

	return reflect.ValueOf(nil), errors.New("could no correctly add document to struct")
}

func deserializeObject(bson []byte, index int, value reflect.Value, isDocument bool) error {
	log.Printf("deserializeObject ======================================== 0")
	// Alright let's parse the fields of the document
	// Decode the length of the buffer
	documentSize := readUInt32(bson, index)

	// Special case of an empty document
	if documentSize == 5 {
		return nil
	}

	// initialIndex
	endIndex := index + int(documentSize)

	// Skip the size document
	index = index + 4

	// Start decoding the fields
	for index < endIndex-1 {
		// Get the bson type
		bsonType := bson[index]

		// Skip bson type
		index = index + 1

		// Read the cstring
		strindex := bytes.IndexByte(bson[index:], 0x00)

		// No 0 byte found error out
		if strindex == -1 {
			return errors.New("could not decode field name, possibly corrupt bson")
		}

		// cast byte array to string
		fieldName := string(bson[index : index+strindex])
		// Adjust index with the string length
		index = index + strindex + 1

		log.Printf("=========== fieldname :: %v type :: %v", fieldName, bsonType)

		// Switch on type to decode
		switch bsonType {
		case byte(bsonString):
			log.Printf("deserializeObject :: string")
			// Read the string size
			stringSize := int(readUInt32(bson, index))
			// Skip string size
			index = index + 4
			// Add to the field value
			err := addValueToFieldStruct(fieldName, value, string(bson[index:index+stringSize-1]), isDocument)
			if err != nil {
				return err
			}
			// Skip last null byte and size of string
			index = index + stringSize
		case byte(bsonDocument):
			log.Printf("deserializeObject :: document")
			// Read the document size
			stringSize := int(readUInt32(bson, index))

			// Add to the field value
			v, err := addDocumentToFieldStruct(fieldName, value, isDocument)
			if err != nil {
				return err
			}

			// Deserialize documents
			err = deserializeObject(bson[index:index+stringSize], 0, v, isDocument)
			if err != nil {
				return err
			}

			// Skip last null byte and size of string
			index = index + stringSize
		case byte(bsonInt32):
			log.Printf("deserializeObject :: int32")
			// Add to the field value
			err := addValueToFieldStruct(fieldName, value, int32(readUInt32(bson, index)), isDocument)
			if err != nil {
				return err
			}

			// Skip the int32 field
			index = index + 4
		}
	}

	// Adjust for last byte
	index = index + 1
	// Return no error
	return nil
	// // Create document node
	// // document := make(map[string]interface{})
	// document := NewDocument()

	// // Special case of an empty document
	// if documentSize == 5 {
	// 	return document, nil
	// }

	// // Skip the size document
	// index = index + 4

	// // Start decoding the fields
	// for index < endIndex-1 {
	// 	// Get the bson type
	// 	bsonType := bson[index]

	// 	// Skip bson type
	// 	index = index + 1

	// 	// Read the cstring
	// 	strindex := bytes.IndexByte(bson[index:], 0x00)

	// 	// No 0 byte found error out
	// 	if strindex == -1 {
	// 		return nil, errors.New("could not decode field name, possibly corrupt bson")
	// 	}

	// 	// cast byte array to string
	// 	fieldName := string(bson[index : index+strindex])
	// 	// Adjust index with the string length
	// 	index = index + strindex + 1

	// 	log.Printf("=========== fieldname :: %v type :: %v", fieldName, bsonType)

	// 	// Switch on type to decode
	// 	switch bsonType {
	// 	default:
	// 		return nil, errors.New(fmt.Sprintf("type [%v] is not a legal bson type", bson[index]))
	// 	case 0x00:
	// 		index = index + 1
	// 		break
	// 	case 0x10:
	// 		// Read the int
	// 		document.Add(fieldName, int32(readUInt32(bson, index)))
	// 		// Skip the int32 field
	// 		index = index + 4
	// 	case 0x09:
	// 		document.Add(fieldName, time.Unix(int64(readUInt64(bson, index)), 0))
	// 		// Skip the int64 field
	// 		index = index + 8
	// 	case 0x08:
	// 		if bson[index] == 0x00 {
	// 			document.Add(fieldName, false)
	// 		} else {
	// 			document.Add(fieldName, true)
	// 		}

	// 		index = index + 1
	// 	case 0x11:
	// 		document.Add(fieldName, &Timestamp{int64(readUInt64(bson, index))})
	// 		// Skip the int64 field
	// 		index = index + 8
	// 	case 0x12:
	// 		document.Add(fieldName, int64(readUInt64(bson, index)))
	// 		// Skip the int64 field
	// 		index = index + 8
	// 	case 0x01:
	// 		document.Add(fieldName, math.Float64frombits(readUInt64(bson, index)))
	// 		// Skip the int64 field
	// 		index = index + 8
	// 	case 0x03:
	// 		// Read the document size
	// 		stringSize := int(readUInt32(bson, index))
	// 		// Deserialize documents
	// 		obj, err := deserializeObject(bson[index:index+stringSize], 0)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		// Set the document
	// 		document.Add(fieldName, obj)
	// 		// Skip last null byte and size of string
	// 		index = index + stringSize
	// 	case 0x04:
	// 		// Read the document size
	// 		stringSize := int(readUInt32(bson, index))
	// 		// Deserialize documents
	// 		obj, err := deserializeArray(bson[index:index+stringSize], 0)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		// Set the document
	// 		document.Add(fieldName, obj)
	// 		// Skip last null byte and size of string
	// 		index = index + stringSize
	// 	case 0x02:
	// 		// Read the string size
	// 		stringSize := int(readUInt32(bson, index))
	// 		// Skip string size
	// 		index = index + 4
	// 		// Add the field
	// 		document.Add(fieldName, string(bson[index:index+stringSize-1]))
	// 		// Skip last null byte and size of string
	// 		index = index + stringSize
	// 	case 0x05:
	// 		// Read the string size
	// 		binarySize := int(readUInt32(bson, index))
	// 		// Skip string size
	// 		index = index + 4
	// 		// Get the subtype
	// 		subStype := bson[index]
	// 		// Skip subtype
	// 		index = index + 1
	// 		// Add the field
	// 		document.Add(fieldName, &Binary{subStype, bson[index : index+binarySize]})
	// 		// Skip last null byte and size of string
	// 		index = index + binarySize
	// 	case 0x07:
	// 		// Add the field
	// 		document.Add(fieldName, &ObjectId{bson[index : index+12]})
	// 		// Skip last null byte and size of string
	// 		index = index + 12
	// 	case 0x06:
	// 		document.Add(fieldName, nil)
	// 	case 0x0a:
	// 		document.Add(fieldName, nil)
	// 	case 0x0d:
	// 		// Read the string size
	// 		stringSize := int(readUInt32(bson, index))
	// 		// Skip string size
	// 		index = index + 4
	// 		// Add the field
	// 		document.Add(fieldName, &Javascript{string(bson[index : index+stringSize-1])})
	// 		// Skip last null byte and size of string
	// 		index = index + stringSize
	// 	case 0x0F:
	// 		// Skip length don't need to decode this
	// 		index = index + 4
	// 		// Read the string size
	// 		stringSize := int(readUInt32(bson, index))
	// 		// Skip string size
	// 		index = index + 4
	// 		// Get the js code string
	// 		jsCode := string(bson[index : index+stringSize-1])
	// 		// Skip last null byte and size of string
	// 		index = index + stringSize
	// 		// Read the document size
	// 		docSize := int(readUInt32(bson, index))
	// 		// Deserialize documents
	// 		obj, err := deserializeObject(bson[index:index+docSize], 0)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		// Create js object
	// 		js := &JavascriptWScope{jsCode, obj}
	// 		// Add the field
	// 		document.Add(fieldName, js)
	// 		// Adjust index
	// 		index = index + docSize
	// 	case 0xff:
	// 		document.Add(fieldName, &Min{})
	// 	case 0x7f:
	// 		document.Add(fieldName, &Max{})
	// 	case 0x0b:
	// 		// Read the cstring
	// 		strindex := bytes.IndexByte(bson[index:], 0x00)

	// 		// No 0 byte found error out
	// 		if strindex == -1 {
	// 			return nil, errors.New("could not decode regexp pattern, possibly corrupt bson")
	// 		}

	// 		// Get regexp pattern
	// 		pattern := string(bson[index : index+strindex])
	// 		// Adjust index with the string length
	// 		index = index + strindex + 1

	// 		// Read the cstring
	// 		strindex = bytes.IndexByte(bson[index:], 0x00)

	// 		// No 0 byte found error out
	// 		if strindex == -1 {
	// 			return nil, errors.New("could not decode regexp options, possibly corrupt bson")
	// 		}

	// 		// Get regexp pattern
	// 		options := string(bson[index : index+strindex])
	// 		// Adjust index with the string length
	// 		index = index + strindex + 1

	// 		// Add regular expression object
	// 		document.Add(fieldName, &RegExp{pattern, options})
	// 	}
	// }

	// // Adjust for last byte
	// index = index + 1
	// // Return the document
	// return document, nil
}

func Deserialize(bson []byte, value reflect.Value) error {
	// Do some basic authentication on the size
	if len(bson) < 5 {
		return errors.New(fmt.Sprintf("Passed in byte slice [%v] is smaller than the minimum size of 5", len(bson)))
	}

	// Decode the length of the buffer
	documentSize := readUInt32(bson, 0)

	// Ensure we have all the bytes
	if documentSize != uint32(len(bson)) {
		return errors.New(fmt.Sprintf("Passed in byte slice [%v] is different in size than encoded bson document length [%v]", len(bson), documentSize))
	}

	// Do we want a pure doc representation instead of serialized struct
	isDocument := false

	log.Printf("======= DERSERIALIZE 0")

	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Switch on the value, validate if we have a document type or not
	switch value.Kind() {
	case reflect.Struct:
		log.Printf("======= DERSERIALIZE 1")
		switch value.Interface().(type) {
		case Document:
			log.Printf("======= DERSERIALIZE 2")
			isDocument = true
		}
	default:
		log.Printf("======= DERSERIALIZE 0:1")
		log.Printf(fmt.Sprintf("cannot deserialize into type %v", value))
		return errors.New(fmt.Sprintf("cannot deserialize into type %v", value))
	}

	// We can start decoding the document
	return deserializeObject(bson, 0, value, isDocument)
}
