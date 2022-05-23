package main

import (
	"log"

	"github.com/nats-io/stan.go"
	// "github.com/nats-io/nats-streaming-server"
)

func main() {
	var clusterID, clientID string
	clusterID = "test-cluster"
	clientID = "clientid"
	
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	defer sc.Close()

	// Simple Synchronous Publisher
	err = sc.Publish("qoo", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	log.Printf("Published async [qoo] : 'msg'\n")
}
