package mongo

import (
	"errors"
	"fmt"
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

func (p *Document) GetString(name string) (string, error) {
	switch elem := p.document[name].(type) {
	default:
		return "", errors.New(fmt.Sprintf("field %v is not a string", name))
	case string:
		return elem, nil
	}
}

func (p *Document) GetInt32(name string) (int32, error) {
	switch elem := p.document[name].(type) {
	default:
		return 0, errors.New(fmt.Sprintf("field %v is not an int32", name))
	case int32:
		return elem, nil
	}
}

func (p *Document) GetDocument(name string) (*Document, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a document", name))
	case *Document:
		return elem, nil
	}
}

func (p *Document) GetBinary(name string) (*Binary, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a binary object", name))
	case *Binary:
		return elem, nil
	}
}

func (p *Document) GetObjectId(name string) (*ObjectId, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a binary object", name))
	case *ObjectId:
		return elem, nil
	}
}

func (p *Document) GetJavascriptField(name string) (*Javascript, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not a javascript object", name))
	case *Javascript:
		return elem, nil
	}
}

func (p *Document) GetJavascriptWScopeField(name string) (*JavascriptWScope, error) {
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

func (p *Document) GetArray(name string) ([]interface{}, error) {
	switch elem := p.document[name].(type) {
	default:
		return nil, errors.New(fmt.Sprintf("field %v is not an array", name))
	case []interface{}:
		return elem, nil
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
