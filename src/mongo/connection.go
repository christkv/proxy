package mongo

import (
	"net"
	"time"
)

type MConnection struct {
	conn net.Conn
}

func NewConnection(address string, timeoutMS int64, ssl bool) (*MConnection, error) {
	// Create a new instance
	conn := new(MConnection)

	// Calculate the timeout
	timeout := time.Duration(timeoutMS) * time.Millisecond

	// Dial a new connection
	socket, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}

	// Assign the socket
	conn.conn = socket
	// Return the connection
	return conn, nil
}

func (p *MConnection) Close() error {
	return p.conn.Close()
}
