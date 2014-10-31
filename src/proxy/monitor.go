package proxy

import (
	"errors"
	"gopkg.in/mgo.v2"
	"log"
	"net"
	"time"
)

// Updated the view of the world based on the mgo
// drivers view of the world
func updateWorldView(context *ConnectionContext, session *mgo.Session, duration time.Duration) error {
	var newMaster = true
	// liveServers := session.LiveServers()
	masterServers := session.LiveMasters()
	liveServers := session.LiveServers()

	if len(masterServers) == 0 {
		return errors.New("no primary found")
	}

	// Validate if our current primary is one of the ones
	// listed in the master list
	for _, mserver := range masterServers {
		if mserver == context.Primary.Address {
			newMaster = false
			break
		}
	}

	//
	// We have a new primary we need to update our entire view of the replicaset
	// 1. Close all connections
	// 2. Reconnect according to the live server list
	if newMaster {
		log.Printf("New primary at %v from %v", masterServers[0], context.Primary.Address)

		context.Primary.Connection.Close()
		// Close all the connection
		for _, con := range context.Secondaries {
			con.Connection.Close()
		}

		// Clean out the context
		context.Primary = nil
		context.Secondaries = make([]*ServerConnection, 0)
		_, err := updateContext(context, session, duration)
		// Setup context
		return err
	}

	// Check if we have any new live servers not covered by the current list
	for _, addr := range liveServers {
		foundServer := false

		// Is the server in the master list then we have a server
		for _, mserver := range masterServers {
			if mserver == addr {
				foundServer = true
				break
			}
		}

		// Is the server in the list of other servers?
		for _, serverConnection := range context.Secondaries {
			if serverConnection.Address == addr {
				foundServer = true
				break
			}
		}

		// No server found
		// Set topology changed, refresh entire set
		if !foundServer {
			log.Printf("Did not locate server %s", addr)
			// Clean out the context
			context.Primary = nil
			context.Secondaries = make([]*ServerConnection, 0)
			_, err := updateContext(context, session, duration)
			return err
		}
	}

	return nil
}

// Open any connections needed
func updateContext(context *ConnectionContext, session *mgo.Session, timeout time.Duration) (*isMasterResult, error) {
	// Timeout duration
	duration := time.Duration(timeout * time.Millisecond)
	// Get the list of live servers
	addresses := session.LiveServers()
	// Is master instance
	isMaster := &isMasterResult{}

	// Establish what server is the master
	err := session.Run("ismaster", isMaster)
	if err != nil {
		log.Fatalf("failed to execute ismaster using mgo %s", err)
		return nil, err
	}

	// For each entry open a tcp connection
	for _, addr := range addresses {
		log.Printf("connect to server %s", addr)
		socket, err := net.DialTimeout("tcp", addr, duration)

		if err != nil {
			log.Printf("failed to connect to server %v", err)
			continue
		}

		// Do we have a primary
		if addr == isMaster.Primary {
			log.Printf("primary found at %s", addr)
			context.Primary = &ServerConnection{addr, socket}
		} else if isSecondary(addr, isMaster.Hosts) {
			context.Secondaries = append(context.Secondaries, &ServerConnection{addr, socket})
		}
	}

	return isMaster, nil
}
