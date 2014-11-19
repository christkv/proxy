package mongo

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"
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
	Scope map[string]interface{}
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

func calculateElementSize(element interface{}) (int, error) {
	size := 0

	// Serialize the document
	switch element := element.(type) {
	default:
		return size, errors.New(fmt.Sprintf("unsupported type %T", element))
	case []interface{}:
		elementSize, err := calculateArraySize(element)
		if err != nil {
			return size, err
		}

		size = size + elementSize
	case map[string]interface{}:
		elementSize, err := calculateObjectSize(element)
		if err != nil {
			return size, err
		}

		size = size + elementSize
	case string:
		size = size + 4 + len(element) + 1
	case int32:
		size = size + 4
	case uint32:
		size = size + 4
	case uint64:
		size = size + 8
	case int64:
		size = size + 8
	case float32:
		size = size + 8
	case float64:
		size = size + 8
	case nil:
		size = size
	case *Binary:
		size = size + 4 + 1 + len(element.Data)
	case bool:
		size = size + 1
	case []byte:
		size = size + 4 + 1 + len(element)
	case *ObjectId:
		size = size + 12
	case *RegExp:
		size = size + len(element.Pattern) + 1 + len(element.Options) + 1
	case *Date:
		size = size + 8
	case *Timestamp:
		size = size + 8
	case time.Time:
		size = size + 8
	case *time.Time:
		size = size + 8
	case *Javascript:
		size = size + len(element.Code) + 4 + 1
	case *JavascriptWScope:
		elementSize, err := calculateObjectSize(element.Scope)
		if err != nil {
			return size, err
		}

		size = size + len(element.Code) + elementSize + 4 + 1 + 4
	case Min:
		size = size
	case Max:
		size = size
	}

	return size, nil
}

func calculateArraySize(array []interface{}) (int, error) {
	size := 5

	// Iterate over all the key values
	for index, value := range array {
		indexStr := strconv.Itoa(index)
		// Add the key size
		size = size + len(indexStr) + 1 + 1
		// Add the size of the actual element
		elementSize, err := calculateElementSize(value)

		if err != nil {
			return size, err
		}

		size = size + elementSize
	}

	return size, nil
}

func calculateObjectSize(document map[string]interface{}) (int, error) {
	size := 5

	// Iterate over all the key values
	for key, value := range document {
		// Add the key size
		size = size + len(key) + 1 + 1
		// Add the size of the actual element
		elementSize, err := calculateElementSize(value)

		if err != nil {
			return size, err
		}

		size = size + elementSize
	}

	return size, nil
}

func packElement(key string, value interface{}, buffer []byte, index int) (int, error) {
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

	// Determine the type
	switch element := value.(type) {
	default:
		return index, errors.New(fmt.Sprintf("unsupported type %T", element))
	case int32:
		log.Printf("int32 serialize")
		buffer[originalIndex] = 0x10
		writeU32(buffer, index+1, uint32(element))
		index = index + 5
	case uint32:
		log.Printf("uint32 serialize")
		buffer[originalIndex] = 0x10
		writeU32(buffer, index+1, element)
		index = index + 5
	case int64:
		log.Printf("int64 serialize")
		buffer[originalIndex] = 0x12
		writeU64(buffer, index+1, uint64(element))
		index = index + 9
	case uint64:
		log.Printf("uint64 serialize")
		buffer[originalIndex] = 0x12
		writeU64(buffer, index+1, element)
		index = index + 9
	case bool:
		log.Printf("bool serialize")
		buffer[originalIndex] = 0x08
		if element {
			buffer[index+1] = 0x01
		} else {
			buffer[index+1] = 0x00
		}
		index = index + 2
	case float32:
		log.Printf("float32 serialize")
		// int64(math.Float64bits(v)
		buffer[originalIndex] = 0x01
		// Get reflection of the value
		reflectType := reflect.ValueOf(element)
		// Convert 32 bit float to string
		floatString := strconv.FormatFloat(reflectType.Float(), 'g', -1, 32)
		// Parse string as 64bit float
		value, _ := strconv.ParseFloat(floatString, 64)
		// Write the float as an uint64
		writeU64(buffer, index+1, math.Float64bits(value))
		index = index + 9
	case float64:
		log.Printf("float64 serialize")
		buffer[originalIndex] = 0x01
		writeU64(buffer, index+1, math.Float64bits(element))
		index = index + 9
	case *RegExp:
		log.Printf("regexp serialize")
		buffer[originalIndex] = 0x0b
		copy(buffer[index+1:], []byte(element.Pattern))
		index = index + 1 + len(element.Pattern)
		buffer[index] = 0x00
		copy(buffer[index+1:], []byte(element.Options))
		index = index + 1 + len(element.Options) + 1
	case *Timestamp:
		log.Printf("timestamp serialize")
		buffer[originalIndex] = 0x11
		writeU64(buffer, index+1, uint64(element.Value))
		index = index + 9
	case *Date:
		log.Printf("date serialize")
		buffer[originalIndex] = 0x09
		writeU64(buffer, index+1, uint64(element.Value))
		index = index + 9
	case *time.Time:
		log.Printf("*time.Time serialize")
		buffer[originalIndex] = 0x09
		writeU64(buffer, index+1, uint64(element.Unix()))
		index = index + 9
	case time.Time:
		log.Printf("time.Time serialize")
		buffer[originalIndex] = 0x09
		writeU64(buffer, index+1, uint64(element.Unix()))
		index = index + 9
	case nil:
		log.Printf("nil serialize")
		buffer[originalIndex] = 0x0a
		index = index + 1
	case string:
		log.Printf("string serialize")
		buffer[originalIndex] = 0x02
		stringBytes := []byte(element)
		writeU32(buffer, index+1, uint32(len(stringBytes)+1))
		copy(buffer[index+5:], stringBytes[:])
		buffer[index+5+len(stringBytes)] = 0x00
		index = index + 4 + len(stringBytes) + 1 + 1
	case []interface{}:
		log.Printf("array serialize %v", index)
		// Set the type of be document
		buffer[originalIndex] = 0x04
		// Get the final values
		in, err := serializeArray(buffer, index+1, element)
		// Serialize the object
		return in + 1, err
	case map[string]interface{}:
		log.Printf("document serialize %v", index)
		// Set the type of be document
		buffer[originalIndex] = 0x03
		// Get the final values
		in, err := serializeObject(buffer, index+1, element)
		// Serialize the object
		return in + 1, err
	case *ObjectId:
		log.Printf("objectid serialize %v", index)
		if len(element.Id) != 12 {
			return 0, errors.New("ObjectId must be a 12 byte array")
		}

		// Set the type of be document
		buffer[originalIndex] = 0x07
		copy(buffer[index+1:], element.Id)
		return index + len(element.Id) + 1, nil
	case []byte:
		log.Printf("[]byte serialize %v", index)
		// Set the type of be document
		buffer[originalIndex] = 0x05
		// Set the size of the binary
		writeU32(buffer, index+1, uint32(len(element)))
		buffer[index+5] = 0x00
		// Write binary
		copy(buffer[index+6:], element[:])
		// Return the length
		return index + len(element) + 5 + 1, nil
	case *Binary:
		log.Printf("binary serialize %v", index)
		// Set the type of be document
		buffer[originalIndex] = 0x05
		// Set the size of the binary
		writeU32(buffer, index+1, uint32(len(element.Data)))
		buffer[index+5] = element.SubType
		// Write binary
		copy(buffer[index+6:], element.Data[:])
		// Return the length
		return index + len(element.Data) + 5 + 1, nil
	case *Javascript:
		log.Printf("javascript no scope serialize %v", index)
		// Set the type of be document
		buffer[originalIndex] = 0x0d
		stringBytes := []byte(element.Code)
		writeU32(buffer, index+1, uint32(len(stringBytes)+1))
		copy(buffer[index+5:], stringBytes[:])
		buffer[index+5+len(stringBytes)] = 0x00
		index = index + 4 + len(stringBytes) + 1 + 1
	case *JavascriptWScope:
		log.Printf("javascript scope serialize %v", index)
		buffer[originalIndex] = 0x0f
		stringBytes := []byte(element.Code)
		// Skip the length
		lengthIndex := index + 1
		index = index + 4

		// Write javascript string
		writeU32(buffer, index+1, uint32(len(stringBytes)+1))
		copy(buffer[index+5:], stringBytes[:])
		buffer[index+5+len(stringBytes)] = 0x00
		index = index + 4 + len(stringBytes) + 1

		// Serialize the scope
		in, err := serializeObject(buffer, index+1, element.Scope)
		index = index + in

		// Write the length
		writeU32(buffer, lengthIndex, uint32(in-lengthIndex+1))
		return in + 1, err
	}

	// Return the index
	return index, nil
}

func serializeObject(buffer []byte, index int, document map[string]interface{}) (int, error) {
	i := index + 4

	// Iterate over all the key values
	for key, value := range document {
		in, err := packElement(key, value, buffer, i)
		if err != nil {
			return i, err
		}
		i = in
	}

	// Final object size
	objectSize := (i - index + 5 - 4)

	// The final object size (the encoded length + terminating 0)
	writeU32(buffer, index, uint32(objectSize))
	// Return no error
	return i, nil
}

func serializeArray(buffer []byte, index int, array []interface{}) (int, error) {
	i := index + 4

	// Iterate over all the key values
	for index, value := range array {
		indexStr := strconv.Itoa(index)
		in, err := packElement(indexStr, value, buffer, i)
		if err != nil {
			return i, err
		}
		i = in
	}

	// Final object size
	objectSize := (i - index + 5 - 4)

	// The final object size (the encoded length + terminating 0)
	writeU32(buffer, index, uint32(objectSize))
	// Return no error
	return i, nil
}

func Serialize(document map[string]interface{}) ([]byte, error) {
	// Calculate the size of the document
	size, err := calculateObjectSize(document)
	if err != nil {
		return nil, err
	}

	log.Printf("size of bson element %v", size)

	// Allocate space
	bson := make([]byte, size)

	// Serialize the object
	_, err = serializeObject(bson, 0, document)
	if err != nil {
		return nil, err
	}

	// Return the bson
	return bson, nil
}
