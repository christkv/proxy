package mongo

import (
	"errors"
	"fmt"
	"reflect"
	// "strings"
)

type encoder struct {
	out       []byte
	index     int
	typeInfos *TypeInfos
}

//
// Initial allocation of byte buffer
const initialAllocationSize = 64

//
// Encode a value into a BSON byte array
//
func (p *BSON) Marshall(doc interface{}, buffer []byte, offset int) ([]byte, error) {
	var out []byte
	var index int

	// Initial allocation of bytes
	if buffer != nil {
		out = buffer
		index = offset
	} else {
		out = make([]byte, initialAllocationSize)
		index = 0
	}

	encoder := &encoder{out, index, p.typeInfos}
	err := encoder.addDoc(reflect.ValueOf(doc))
	if err != nil {
		return nil, err
	}

	// Return the correct slice with the serialized result
	return encoder.out[offset:encoder.index], nil
}

func (p *encoder) packString(originalIndex int, value reflect.Value) int {
	// Get the string
	str := value.String()
	// Set the type
	p.out[originalIndex] = 0x02
	// Get the string bytes
	stringBytes := []byte(str)
	// Write the string length
	writeU32(p.out, p.index+1, uint32(len(stringBytes)+1))
	// Write bytes with bounds checking
	p.writeBytes(stringBytes[:], p.index+5)
	// Set end 0 byte
	p.out[p.index+5+len(stringBytes)] = 0x00
	// Return new index position
	return p.index + 4 + len(stringBytes) + 1 + 1
}

func (p *encoder) writeBytes(bytes []byte, index int) {
	// We need to allocate more memory
	if len(p.out)-index < len(bytes) {
		// Allocate a new buffer
		memory := make([]byte, len(bytes)+initialAllocationSize+len(p.out))
		// Copy existing buffer into it
		copy(memory[0:], p.out[0:index])
		// Point to new buffer
		p.out = memory
	}

	// Write the bytes into the buffer
	copy(p.out[index:], bytes[:])
}

func (p *encoder) packElement(key string, value reflect.Value) error {
	strbytes := []byte(key)
	// Save a pointer to the first byte index
	originalIndex := p.index
	// Skip type
	p.index = p.index + 1

	// Write bytes with bounds checking
	p.writeBytes(strbytes[0:], p.index)
	// Null terminate the string
	p.out[p.index+len(strbytes)] = 0

	// Update the index with the field length
	p.index = p.index + len(strbytes)

	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Reflect on the type
	switch value.Kind() {
	case reflect.String:
		p.index = p.packString(originalIndex, value)
	case reflect.Int32:
		p.out[originalIndex] = 0x10
		writeU32(p.out, p.index+1, uint32(value.Int()))
		p.index = p.index + 5
	case reflect.Struct:
		switch value.Interface().(type) {
		case Document:
			// Set the type of be document
			p.out[originalIndex] = 0x03
			// Skip initial byte
			p.index = p.index + 1
			// Get the final values
			return p.addDoc(value)

		default:
			// Set the type of be document
			p.out[originalIndex] = 0x03
			// Skip initial byte
			p.index = p.index + 1
			// Get the final values
			return p.addDoc(value)
		}
	default:
		return errors.New(fmt.Sprintf("could not recognize the type %v", value.Kind()))
	}

	return nil
}

func (p *encoder) addDoc(value reflect.Value) error {
	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Switch on the value
	switch value.Kind() {
	case reflect.Struct:
		// Save current index for the writing of the total size of the doc
		originalIndex := p.index

		// Skip the 4 for size bytes
		p.index = p.index + 4

		// Check if we have the Document type or a normal struct
		switch doc := value.Interface().(type) {
		case Document:
			// Iterate over all the key values
			for _, key := range doc.fields {
				// Get the value
				fieldValue := doc.document[key]
				// Add the size of the actual element
				err := p.packElement(key, reflect.ValueOf(fieldValue))
				if err != nil {
					return err
				}
			}
		default:
			typeInfo := parseTypeInformation(p.typeInfos, value)
			// Let's iterate over all the fields
			for j := 0; j < typeInfo.NumberOfField; j++ {
				// Get the field value
				fieldValue := value.Field(j)
				// Get field type
				fieldType := typeInfo.FieldsByIndex[j]
				// Get the field name
				key := fieldType.MetaDataName
				// Add the size of the actual element
				err := p.packElement(key, fieldValue)
				if err != nil {
					return err
				}
			}
		}

		// Skip last null byte
		p.index = p.index + 1
		// Write the totalSize of the document
		writeU32(p.out, originalIndex, uint32(p.index-originalIndex))
	default:
		return errors.New(fmt.Sprintf("BSON struct type %T not supported for serialization", value))
	}

	return nil
}
