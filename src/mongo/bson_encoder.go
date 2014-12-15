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
	// Copy the string bytes to the bson buffer
	copy(p.out[p.index+5:], stringBytes[:])
	// Set end 0 byte
	p.out[p.index+5+len(stringBytes)] = 0x00
	// Return new index position
	return p.index + 4 + len(stringBytes) + 1 + 1
}

func (p *encoder) packElement(key string, value reflect.Value) error {
	// fmt.Printf("packElement %v with value %v of kind %v\n", key, value, value.Kind())
	strbytes := []byte(key)
	// Save a pointer to the first byte index
	originalIndex := p.index
	// Skip type
	p.index = p.index + 1
	// Copy the string into the buffer
	copy(p.out[p.index:], strbytes[0:])
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
		// fmt.Printf("packElement ================ got string %v = %v\n", key, value.String())
		p.index = p.packString(originalIndex, value)
	case reflect.Int32:
		// fmt.Printf("packElement ================ got int32 %v = %v\n", key, value.Int())
		p.out[originalIndex] = 0x10
		writeU32(p.out, p.index+1, uint32(value.Int()))
		p.index = p.index + 5
	case reflect.Struct:
		// fmt.Printf("packElement ================ got document %v\n", value.Interface())
		switch value.Interface().(type) {
		case Document:
			// fallthrough
			// fmt.Printf("packElement ================ got document\n")
			// Set the type of be document
			p.out[originalIndex] = 0x03
			// Skip initial byte
			p.index = p.index + 1
			// Get the final values
			return p.addDoc(value)

		default:
			// fmt.Printf("packElement ================ got struct %v = %v at %v\n", key, value, p.index)
			// Set the type of be document
			p.out[originalIndex] = 0x03
			// Skip initial byte
			p.index = p.index + 1
			// Get the final values
			return p.addDoc(value)
		}
	default:
		// fmt.Printf("could not recognize the type %v\n", value.Kind())
		return errors.New(fmt.Sprintf("could not recognize the type %v", value.Kind()))
	}

	return nil
}

func (p *encoder) addDoc(value reflect.Value) error {
	// fmt.Printf("============================== serialize addDoc = %v\n", value.Kind())
	// We have a pointer get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Switch on the value
	switch value.Kind() {
	case reflect.Struct:
		// fmt.Printf("============================== serialize addDoc struct\n")
		// Save current index for the writing of the total size of the doc
		originalIndex := p.index

		// Skip the 4 for size bytes
		p.index = p.index + 4

		// Check if we have the Document type or a normal struct
		switch doc := value.Interface().(type) {
		case Document:
			// fmt.Printf("============================== serialize document\n")

			// Iterate over all the key values
			for _, key := range doc.fields {
				// fmt.Printf("key :: %v", key)
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

			// fmt.Printf("============================== serialize interface\n")
			numberOfField := value.NumField()

			// Let's iterate over all the fields
			for j := 0; j < numberOfField; j++ {
				// Get the field value
				fieldValue := value.Field(j)
				// Get field type
				fieldType := typeInfo.FieldsByIndex[j]
				// Get the field name
				key := fieldType.Name

				// // Get the tag
				// tag := fieldType.Tag.Get("bson")
				// // Split the tag into parts
				// parts := strings.Split(tag, ",")

				// // Override the key if the metadata has one
				// if len(parts) > 0 && parts[0] != "" {
				// 	key = parts[0]
				// }

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
