package mongo

import (
	"errors"
	"fmt"
	"time"
)

type KeyValue struct {
	Name  string
	Value interface{}
}

func NewDocument() *Document {
	document := &Document{make([]string, 0), make(map[string]interface{})}
	return document
}

type Document struct {
	fields   []string
	document map[string]interface{}
}

func (p *Document) FieldCount() int {
	return len(p.fields)
}

func (p *Document) Add(name string, value interface{}) {
	p.fields = append(p.fields, name)
	p.document[name] = value
}

func (p *Document) FieldAsFloat64(name string) (float64, error) {
	switch elem := p.document[name].(type) {
	default:
		return 0, errors.New(fmt.Sprintf("field %v is not a float64", name))
	case float64:
		return elem, nil
	}
}

func (p *Document) FieldAsFloat32(name string) (float32, error) {
	switch elem := p.document[name].(type) {
	default:
		return 0, errors.New(fmt.Sprintf("field %v is not a float32", name))
	case float64:
		return float32(elem), nil
	}
}

func (p *Document) FieldAsRegExp(name string) (*RegExp, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a regexp", name))
	case *RegExp:
		return elem, nil
	}
}

func (p *Document) FieldAsBool(name string) (bool, error) {
	switch elem := p.document[name].(type) {
	default:
		return false, errors.New(fmt.Sprintf("field %v is not a bool", name))
	case bool:
		return elem, nil
	}
}

func (p *Document) FieldAsNil(name string) (interface{}, error) {
	switch elem := p.document[name].(type) {
	default:
		return false, errors.New(fmt.Sprintf("field %v is not nil", name))
	case nil:
		return elem, nil
	}
}

func (p *Document) FieldAsString(name string) (string, error) {
	switch elem := p.document[name].(type) {
	default:
		return "", errors.New(fmt.Sprintf("field %v is not a string", name))
	case string:
		return elem, nil
	}
}

func (p *Document) FieldAsInt32(name string) (int32, error) {
	switch elem := p.document[name].(type) {
	default:
		return 0, errors.New(fmt.Sprintf("field %v is not an int32", name))
	case int32:
		return elem, nil
	}
}

func (p *Document) FieldAsDocument(name string) (*Document, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a document", name))
	case *Document:
		return elem, nil
	}
}

func (p *Document) FieldAsBinary(name string) (*Binary, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a binary object", name))
	case *Binary:
		return elem, nil
	}
}

func (p *Document) FieldAsObjectId(name string) (*ObjectId, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a binary object", name))
	case *ObjectId:
		return elem, nil
	}
}

func (p *Document) FieldAsJavascript(name string) (*Javascript, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a javascript object", name))
	case *Javascript:
		return elem, nil
	}
}

func (p *Document) FieldAsJavascriptWScope(name string) (*JavascriptWScope, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a javascript object", name))
	case *JavascriptWScope:
		return elem, nil
	}
}

func (p *Document) FieldAsMin(name string) (*Min, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a min object", name))
	case *Min:
		return elem, nil
	}
}

func (p *Document) FieldAsMax(name string) (*Max, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a max object", name))
	case *Max:
		return elem, nil
	}
}

func (p *Document) FieldAsArray(name string) ([]interface{}, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not an array", name))
	case []interface{}:
		return elem, nil
	}
}

func (p *Document) FieldAsTime(name string) (time.Time, error) {
	switch elem := p.document[name].(type) {
	default:
		return time.Unix(0, 0), errors.New(fmt.Sprintf("field %v is not a time instance", name))
	case time.Time:
		return elem, nil
	}
}

func (p *Document) FieldAsTimestamp(name string) (*Timestamp, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a timestamp", name))
	case *Timestamp:
		return elem, nil
	}
}

func (p *Document) FieldAsInt64(name string) (int64, error) {
	switch elem := p.document[name].(type) {
	default:
		return 0, errors.New(fmt.Sprintf("field %v is not an int64", name))
	case int64:
		return elem, nil
	case uint64:
		return int64(elem), nil
	}
}

func (p *Document) FieldAsUInt64(name string) (uint64, error) {
	switch elem := p.document[name].(type) {
	default:
		return 0, errors.New(fmt.Sprintf("field %v is not an uint64", name))
	case uint64:
		return elem, nil
	case int64:
		return uint64(elem), nil
	}
}

func (p *Document) FieldsInOrder() []KeyValue {
	// Create an array of key values to return
	keyValues := make([]KeyValue, len(p.fields))

	// Convert all fields into key values
	for i, key := range p.fields {
		keyValues[i] = KeyValue{key, p.document[key]}
	}

	// Return all the key values
	return keyValues
}
