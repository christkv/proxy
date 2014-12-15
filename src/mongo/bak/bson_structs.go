package bak

// import (
// "errors"
// "log"
// "reflect"
// "strings"
// )

// func SerializeStruct(obj interface{}, bson []byte, offset int) ([]byte, error) {
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
// 			// Calculate the size of the document
// 			size, err := CalculateObjectSize(typeof)
// 			if err != nil {
// 				return nil, err
// 			}

// 			log.Printf("size of bson element %v", size)

// 			// Allocate space
// 			bson = make([]byte, size)
// 		}

// 		numberOfField := typeof.NumField()

// 		log.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\nnumber of fields off struct %v", numberOfField)

// 		// Let's iterate over all the fields
// 		for i := 0; i < numberOfField; i++ {
// 			// Get the field value
// 			fieldValue := document.Field(i)
// 			fieldType := document.Type().Field(i)
// 			// Get the field name
// 			key := fieldType.Name
// 			// Get the tag
// 			tag := fieldType.Tag.Get("bson")
// 			// Split the tag into parts
// 			parts := strings.Split(tag, ",")

// 			// Override the key if the metadata has one
// 			if len(parts) > 0 && parts[0] != "" {
// 				key = parts[0]
// 			}

// 			log.Printf("calculate size for field %v of type %v with tag %v", key, fieldValue, tag)

// 			// Add the length of the name of the field
// 			size = size + len(key) + 1 + 1

// 			// Add the size of the actual element
// 			elementSize, err := calculateElementSize(fieldValue)

// 			if err != nil {
// 				return size, err
// 			}

// 			size = size + elementSize
// 		}

// 		// // Using the BSON serialize the structure into the buffer
// 		// _, err := serializeObject(bson[offset:], 0, obj)
// 		// if err != nil {
// 		// 	return nil, err
// 		// }

// 		// // Return the bson
// 		// return bson, nil

// 		// // Calculate the size of the document
// 		// size, err := CalculateObjectSize(typeof)
// 		// if err != nil {
// 		// 	return nil, err
// 		// }

// 		// log.Printf("################## size %v", size)

// 		// // Allocate space for the bson document
// 		// bson := make([]byte, size)

// 		// Write the size in
// 		// writeU32(bson, 0, uint32(size))

// 		// Return bson
// 		return bson, nil
// 	default:
// 		return nil, errors.New("Unsupported Serialization Mode")
// 	}
// }
