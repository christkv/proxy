package proxy

import (
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
	IsMaster       bool
	Secondary      bool
	Primary        string
	Hosts          []string
	Passives       []string
	Tags           bson.D
	Msg            string
	MaxWireVersion int `bson:"maxWireVersion"`
}

type ConnectionContext struct {
	Primary     net.Conn
	Secondaries []net.Conn
}

func isSecondary(addr string, addresses []string) bool {
	for _, a := range addresses {
		if a == addr {
			return true
		}
	}

	return false
}

func HandleConnection(set *ReplSet, conn net.Conn) {
	// Set up socket connections
	addresses := set.Session.LiveServers()

	// Timeout duration
	duration := time.Duration(set.Timeout * time.Millisecond)

	// Is master instance
	isMaster := &isMasterResult{}

	// Clean up connection on exit
	defer conn.Close()

	// Establish what server is the master
	err := set.Session.Run("ismaster", isMaster)
	if err != nil {
		log.Fatalf("failed to execute ismaster using mgo %s", err)
		return
	}

	// Create connection context
	context := &ConnectionContext{}
	context.Secondaries = make([]net.Conn, 0)

	// For each entry open a tcp connection
	for _, addr := range addresses {
		socket, err := net.DialTimeout("tcp", addr, duration)

		if err != nil {
			log.Printf("failed to connect to server %v", err)
			continue
		}

		// Do we have a primary
		if addr == isMaster.Primary {
			log.Printf("primary found at %s", addr)
			context.Primary = socket
		} else if isSecondary(addr, isMaster.Hosts) {
			context.Secondaries = append(context.Secondaries, socket)
		}
	}

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

		// Let's unpack the wire message header
		opCode := readInt32(wireMessage[8:12])

		// Cluster connection
		connection := context.Primary

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
