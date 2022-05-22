package main

import (
	"log"

	"github.com/nats-io/stan.go"
	// "fmt"
)

func main() {
	var clusterID, serverID string
	clusterID = "test-cluster"
	serverID = ""
	
	sc, err := stan.Connect(clusterID, serverID)
	if err != nil {
		log.Fatalf("Can't connect publish %s", err.Error())
	}

	// Simple Synchronous Publisher
	sc.Publish("foo", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming
	
	// // Simple Async Subscriber
	// sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
	// })
	
	// // Unsubscribe
	// sub.Unsubscribe()
	
	// Close connection
	sc.Close()
}