package mongo

import (
	"errors"
	"fmt"
	"log"
)

type ObjectId struct {
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

func writeU32(buffer []byte, index int, size uint32) {
	buffer[index+3] = byte((size >> 24) & 0xff)
	buffer[index+2] = byte((size >> 16) & 0xff)
	buffer[index+1] = byte((size >> 8) & 0xff)
	buffer[index] = byte(size & 0xff)
}

func calculateElementSize(element interface{}) (int, error) {
	size := 0

	// Serialize the document
	switch element := element.(type) {
	default:
		return size, errors.New(fmt.Sprintf("unsupported type %T", element))
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
		size = size + 4
	case int64:
		size = size + 4
	case Binary:
		size = size + 4 + 1 + len(element.Data)
	case ObjectId:
		size = size + 12
	case Javascript:
		size = size + len(element.Code)
	case JavascriptWScope:
		elementSize, err := calculateObjectSize(element.Scope)
		if err != nil {
			return size, err
		}

		size = size + len(element.Code) + elementSize
	case Min:
		size = size
	case Max:
		size = size
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
	case string:
		log.Printf("string serialize")
		buffer[originalIndex] = 0x02
		stringBytes := []byte(element)
		writeU32(buffer, index+1, uint32(len(stringBytes)+1))
		copy(buffer[index+5:], stringBytes[:])
		buffer[index+5+len(stringBytes)] = 0x00
		index = index + 4 + len(stringBytes) + 1 + 1
	case map[string]interface{}:
		log.Printf("document serialize %v", index)
		// Set the type of be document
		buffer[originalIndex] = 0x03
		// Get the final values
		in, err := serializeObject(buffer, index+1, element)
		// Serialize the object
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
