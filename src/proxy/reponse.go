package proxy

import (
	"gopkg.in/mgo.v2/bson"
)

func CreateResponseMessage(requestId []byte, obj interface{}) ([]byte, error) {
	// Serialize to bson
	data, err := bson.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Total length
	msgLength := len(data) + 16 + 20

	// Create the reponse message
	ismasterCommandBytes := make([]byte, 0)
	// 16 byte header
	ismasterCommandBytes = addInt32(ismasterCommandBytes, int32(msgLength))
	ismasterCommandBytes = append(ismasterCommandBytes, []byte{0, 0, 0, 0}...)
	ismasterCommandBytes = append(ismasterCommandBytes, requestId...)
	ismasterCommandBytes = addInt32(ismasterCommandBytes, int32(OP_REPLY))
	// OP REPLY FIELDS
	ismasterCommandBytes = append(ismasterCommandBytes, []byte{0, 0, 0, 0}...)
	ismasterCommandBytes = append(ismasterCommandBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
	ismasterCommandBytes = append(ismasterCommandBytes, []byte{0, 0, 0, 0}...)
	ismasterCommandBytes = addInt32(ismasterCommandBytes, int32(1))
	ismasterCommandBytes = append(ismasterCommandBytes, data...)
	return ismasterCommandBytes, nil
}
