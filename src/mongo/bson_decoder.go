package mongo

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func (p *BSON) Unmarshal(bson []byte, obj interface{}) error {
	// Do some basic authentication on the size
	if len(bson) < 5 {
		return errors.New(fmt.Sprintf("Passed in byte slice [%v] is smaller than the minimum size of 5", len(bson)))
	}

	// Get value type
	value := reflect.ValueOf(obj)

	// Decode the length of the buffer
	documentSize := readUInt32(bson, 0)

	// Ensure we have all the bytes
	if documentSize != uint32(len(bson)) {
		return errors.New(fmt.Sprintf("Passed in byte slice [%v] is different in size than encoded bson document length [%v]", len(bson), documentSize))
	}

	// Do we want a pure doc representation instead of serialized struct
	isDocument := false

	// Check if we have a *Document or Document instance
	switch value.Interface().(type) {
	case *Document:
		isDocument = true
	}

	// If we have a pointer get to the value object
	if isDocument == false && value.Kind() == reflect.Ptr {
		// v := reflect.ValueOf(out)
		switch value.Kind() {
		case reflect.Ptr:
			//  fallthrough
			// case reflect.Map:
			value = value.Elem()
		case reflect.Struct:
			return errors.New("must be a pointer to a struct or Document")
		default:
			return errors.New("must be a pointer to a struct or Document")
		}
	}

	return p.deserializeObject(bson, 0, value, isDocument)
}

func (p *BSON) deserializeObject(bson []byte, index int, value reflect.Value, isDocument bool) error {
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

		// Switch on type to decode
		switch bsonType {
		case byte(bsonString):
			// Read the string size
			stringSize := int(readUInt32(bson, index))
			// Skip string size
			index = index + 4
			// Add to the field value
			err := p.addValueToFieldStruct(fieldName, value, string(bson[index:index+stringSize-1]), isDocument)
			if err != nil {
				return err
			}
			// Skip last null byte and size of string
			index = index + stringSize
		case byte(bsonDocument):
			// Read the document size
			stringSize := int(readUInt32(bson, index))

			// Add to the field value
			v, err := p.addDocumentToFieldStruct(fieldName, value, isDocument)
			if err != nil {
				return err
			}

			// Deserialize documents
			err = p.deserializeObject(bson[index:index+stringSize], 0, v, isDocument)
			if err != nil {
				return err
			}

			// Skip last null byte and size of string
			index = index + stringSize
		case byte(bsonInt32):
			// Add to the field value
			err := p.addValueToFieldStruct(fieldName, value, int32(readUInt32(bson, index)), isDocument)
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
}

func (p *BSON) addValueToFieldStruct(fieldName string, obj reflect.Value, value interface{}, isDocument bool) error {
	if isDocument {
		switch t := obj.Interface().(type) {
		case *Document:
			t.Add(fieldName, value)
		}
	} else {
		// Get the type info
		typeInfo := parseTypeInformation(p.typeInfos, obj)
		structFieldName := typeInfo.Fields[fieldName].Name

		if obj.Kind() == reflect.Ptr {
			obj = obj.Elem()
		}

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

func (p *BSON) addDocumentToFieldStruct(fieldName string, obj reflect.Value, isDocument bool) (reflect.Value, error) {
	if isDocument {
		switch t := obj.Interface().(type) {
		case *Document:
			// Create a new document
			doc := NewDocument()
			// Add the field
			t.Add(fieldName, doc)
			// Return the new value
			return reflect.ValueOf(doc), nil
		}
	} else {
		// Get the type info
		typeInfo := parseTypeInformation(p.typeInfos, obj)
		structFieldName := typeInfo.Fields[fieldName].Name

		// Set the field value on the struct (just set it hard)
		field := obj.FieldByName(structFieldName)
		fieldType, _ := obj.Type().FieldByName(structFieldName)

		// We did not find the field on the struct
		if field.Kind() == reflect.Invalid {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("field %v not found on struct %v", fieldName, obj))
		}

		switch field.Kind() {
		case reflect.Ptr:
			// Get raw type
			underlyingType := fieldType.Type.Elem()

			// Create an instance
			value := reflect.New(underlyingType)

			// Set the field on the struct
			field.Set(value)

			// Return the new value
			return value, nil
		}
	}

	return reflect.ValueOf(nil), errors.New("could no correctly add document to struct")
}
