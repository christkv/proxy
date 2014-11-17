package mongo

import (
	"testing"
)

func TestSuccessfulConnection(t *testing.T) {
	conn, err := NewConnection("localhost:27017", 1000, true)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB")
	}

	conn.Close()
}

func TestShouldSuccessfullyExecuteCommand(t *testing.T) {
	conn, err := NewConnection("localhost:27017", 1000, true)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB")
	}

	// Encode ismaster command

	conn.Close()
}
