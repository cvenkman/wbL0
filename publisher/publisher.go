package main

import (
	"log"

	"github.com/nats-io/stan.go"
	"flag"
)

// var usageStr = `
// Usage: publisher [options]

// Options:
// 	-cl		<cluster name>	NATS Streaming cluster name
// 	-id		<client ID>		NATS Streaming client ID
// 	-ch		<channel>		The NATS Streaming channel to create
// 	-msg	<message>		NATS Streaming cluster name
// `

func main() {
	var clusterID, clientID, channel, msg string
	flag.StringVar(&clusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-publisher", "The NATS Streaming client ID to connect with")
	flag.StringVar(&channel, "ch", "test-channel", "The NATS Streaming channel to create")
	flag.StringVar(&msg, "msg", "hello!", "Data to publish")
	flag.Parse()

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	defer sc.Close()

	err = sc.Publish(channel, []byte(msg)) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	log.Printf("Published [%s] : '%s'\n", channel, msg)
}
