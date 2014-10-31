package proxy

import (
	"bytes"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"time"
)

const OP_REPLY = 1
const OP_MSG = 1000
const OP_UPDATE = 2001
const OP_INSERT = 2002
const OP_QUERY = 2004
const OP_GET_MORE = 2005
const OP_DELETE = 2006
const OP_KILL_CURSORS = 2007

type isMasterResult struct {
	Ok                  int
	IsMaster            bool
	Secondary           bool     `bson:",omitempty"`
	Primary             string   `bson:",omitempty"`
	Hosts               []string `bson:",omitempty"`
	Passives            []string `bson:",omitempty"`
	Tags                bson.D   `bson:",omitempty"`
	Msg                 string   `bson:",omitempty"`
	MaxMessageSizeBytes int
	MaxWireVersion      int `bson:"maxWireVersion"`
	MaxBsonObjectSize   int
	LocalTime           time.Time
}

type readPreference struct {
	Mode string
	Tags []map[string]string
}

type ServerConnection struct {
	Address    string
	Connection net.Conn
}

type ConnectionContext struct {
	Primary     *ServerConnection
	Secondaries []*ServerConnection
}

func isSecondary(addr string, addresses []string) bool {
	for _, a := range addresses {
		if a == addr {
			return true
		}
	}

	return false
}

func addInt32(b []byte, i int32) []byte {
	return append(b, byte(i), byte(i>>8), byte(i>>16), byte(i>>24))
}

func addInt64(b []byte, i int64) []byte {
	return append(b, byte(i), byte(i>>8), byte(i>>16), byte(i>>24),
		byte(i>>32), byte(i>>40), byte(i>>48), byte(i>>56))
}

func HandleConnection(set *ReplSet, conn net.Conn) {
	var isMasterBytes = []byte("isMaster")
	var ismasterBytes = []byte("ismaster")
	var readPreferenceBytes = []byte("$readPreference")
	var readPreferencebytes = []byte("$readpreference")
	var connection net.Conn

	// Set up socket connections
	addresses := set.Session.LiveServers()
	masters := set.Session.LiveMasters()

	log.Printf("live servers [%v], Live masters[%v]", addresses, masters)

	// // Timeout duration
	// duration := time.Duration(set.Timeout * time.Millisecond)

	// // Is master instance
	// isMaster := &isMasterResult{}

	// Clean up connection on exit
	defer conn.Close()

	// // Establish what server is the master
	// err := set.Session.Run("ismaster", isMaster)
	// if err != nil {
	// 	log.Fatalf("failed to execute ismaster using mgo %s", err)
	// 	return
	// }

	// Create connection context
	context := &ConnectionContext{}
	context.Secondaries = make([]*ServerConnection, 0)

	isMaster, err := updateContext(context, set.Session, set.Timeout)
	if err != nil {
		log.Fatalf("failed to update the world %s", err)
		return
	}

	// // For each entry open a tcp connection
	// for _, addr := range addresses {
	// 	socket, err := net.DialTimeout("tcp", addr, duration)

	// 	if err != nil {
	// 		log.Printf("failed to connect to server %v", err)
	// 		continue
	// 	}

	// 	// Do we have a primary
	// 	if addr == isMaster.Primary {
	// 		log.Printf("primary found at %s", addr)
	// 		context.Primary = &ServerConnection{addr, socket}
	// 	} else if isSecondary(addr, isMaster.Hosts) {
	// 		context.Secondaries = append(context.Secondaries, &ServerConnection{addr, socket})
	// 	}
	// }

	// Start reading of messages
	for {
		messageSizeBytes := make([]byte, 4)
		n, err := conn.Read(messageSizeBytes)

		// We have an error, close socket and return
		if err != nil || int32(n) != 4 {
			log.Printf("failed to read enough bytes to establish message size from connection %v", err)
			break
		}

		// Get the message size
		messageSize := readInt32(messageSizeBytes)

		// We know the size of the message, read the entire message into memory
		wireMessage := make([]byte, messageSize-4)
		n, err = conn.Read(wireMessage)

		// We had an error during the reading of the message, close socket and return
		if err != nil || int32(n) != (messageSize-4) {
			log.Printf("failed to read enough bytes for wire protocol message from connection %v", err)
			break
		}

		log.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		// Update our view of the world to match the one from the mgo driver
		err = updateWorldView(context, set.Session, set.Timeout)
		if err != nil {
			log.Printf("failed to update view of the world %v", err)
			break
		}

		// Let's unpack the wire message header
		opCode := readInt32(wireMessage[8:12])
		// Get possible indexes
		index := -1
		index1 := bytes.Index(wireMessage, readPreferenceBytes)
		index2 := bytes.Index(wireMessage, readPreferencebytes)
		if index1 != -1 {
			index = index1
		} else if index2 != -1 {
			index = index2
		}

		// Default to primary
		connection = context.Primary.Connection

		// Look for readPreference provided by client in the message
		if index != -1 {
			log.Printf("server requesting read preference from replicaset")
			// We know that infront of a bson object (which $readPreference is)
			// there is the field type 1 byte
			// and the total object length 4 bytes
			// let's get the bson object
			readPreferenceSize := readInt32(wireMessage[index+len(readPreferenceBytes)+1:])
			// Get the readpreference bson doc and deserialize it
			readPrefBytes := wireMessage[index+len(readPreferenceBytes)+1 : index+len(readPreferenceBytes)+1+int(readPreferenceSize)]
			// Read Preference
			readPref := &readPreference{}
			// Unmarshal
			err = bson.Unmarshal(readPrefBytes, readPref)
			if err != nil {
				log.Printf("failed to deserialize the readPreference object %v", err)
				return
			}

			// Print the read preference
			// log.Printf("readPreference : %+v", readPref)
			// If we have secondary read preference
			if (readPref.Mode == "secondary" ||
				readPref.Mode == "secondaryPreferred" ||
				readPref.Mode == "nearest") && len(context.Secondaries) > 0 {
				log.Printf("execute operation against secondary")
				connection = context.Secondaries[0].Connection
			}
		}

		// Determine if this is the ismaster command from a driver
		if bytes.Index(wireMessage, isMasterBytes) != -1 || bytes.Index(wireMessage, ismasterBytes) != -1 {

			// Header fields
			requestId := wireMessage[0:4]

			// Create command
			var ismasterCmd = &isMasterResult{
				Ok:                  1,
				IsMaster:            true,
				MaxMessageSizeBytes: isMaster.MaxMessageSizeBytes,
				MaxBsonObjectSize:   isMaster.MaxBsonObjectSize,
				Msg:                 "isdbgrid",
				MaxWireVersion:      isMaster.MaxWireVersion,
				LocalTime:           isMaster.LocalTime,
			}

			ismasterCommandBytes, err := CreateResponseMessage(requestId, ismasterCmd)
			if err != nil {
				log.Printf("failed to create ismaster command %v", err)
				break
			}

			// Write ismaster response
			conn.Write(ismasterCommandBytes)
			continue
		}

		// If it's write commands we need to direct it to the primary
		if opCode == OP_INSERT || opCode == OP_UPDATE || opCode == OP_DELETE || opCode == OP_KILL_CURSORS {
			connection.Write(messageSizeBytes)
			connection.Write(wireMessage)
		} else if opCode == OP_GET_MORE || opCode == OP_QUERY {
			connection.Write(messageSizeBytes)
			connection.Write(wireMessage)

			// Read the response from the connection
			responseMessageSizeBytes := make([]byte, 4)
			n, err = connection.Read(responseMessageSizeBytes)

			// We have an error, close socket and return
			if err != nil || int32(n) != 4 {
				log.Printf("failed to read enough bytes to establish message size from connection %v", err)
				break
			}

			responseMessageSize := readInt32(responseMessageSizeBytes)
			responseMessageBytes := make([]byte, responseMessageSize-4)

			// Read the rest of the response
			n, err = connection.Read(responseMessageBytes)

			// We have an error, close socket and return
			if err != nil || int32(n) != (responseMessageSize-4) {
				log.Printf("failed to read enough bytes to establish message size from connection %v", err)
				break
			}

			// Write message to initial connection
			conn.Write(responseMessageSizeBytes)
			conn.Write(responseMessageBytes)
		} else {
			log.Fatalf("opcode %v not supported", opCode)
			break
		}
	}
}

func readInt32(b []byte) int32 {
	return int32((uint32(b[0]) << 0) |
		(uint32(b[1]) << 8) |
		(uint32(b[2]) << 16) |
		(uint32(b[3]) << 24))
}
