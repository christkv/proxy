package main

import (
	// "fmt"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"time"
	// "os"
	"proxy"
	// "strings"
)

func main() {
	// Parser flags
	var uri string
	var timeout int

	// Proxy command
	var proxyCmd = &cobra.Command{
		Use:   "mongor",
		Short: "Mongor is a replicaset proxy",
		Long: `Mongor is a replicaset proxy that allows for simple
            connectivity to a replicaset for drivers wanting to support
            replicaset connectivity without implementing the complex
            monitoring or who need to centralize connections due to single
            threaded application platforms`,
		Run: func(cmd *cobra.Command, args []string) {
			// Listen to the tcp socket
			ln, err := net.Listen("tcp", ":50000")
			if err != nil {
				log.Fatalf("%s", err)
			}

			// Create ReplSet
			set := proxy.NewReplSet(uri, time.Duration(timeout))
			// Attempt to Connect to the replicaset
			err = set.Start()
			if err != nil {
				log.Fatalf("failed to connect to replicaset %s", err)
			}

			// Accept incoming socket connection
			for {
				conn, err := ln.Accept()
				if err != nil {
					log.Fatalf("%s", err)
				}

				// Fire off our handler
				go proxy.HandleConnection(set, conn)
			}
		},
	}

	// Set default timeout
	timeout = 10000

	// Set up the uri flag
	proxyCmd.Flags().StringVarP(&uri, "uri", "u", "mongodb://localhost:31000/admin?maxPoolSize=1", "replicaset connection uri")
	proxyCmd.Execute()
}
