package main

import (
	"log"

	"github.com/nats-io/stan.go"
	"flag"
	"io/ioutil"
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
	var clusterID, clientID, channel, dataPath string
	flag.StringVar(&clusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-publisher", "The NATS Streaming client ID to connect with")
	flag.StringVar(&channel, "ch", "test-channel", "The NATS Streaming channel to create")
	// flag.StringVar(&msg, "msg", "hello!", "Data to publish")
	flag.StringVar(&dataPath, "data", "model/model.json", "Json data(model) to publish")
	flag.Parse()

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect publish: %s", err.Error())
	}
	defer sc.Close()

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't read json model: %s", err.Error())
	}

	err = sc.Publish(channel, data) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect publish: %s", err.Error())
	}
	log.Printf("Published [%s]: '%s'\n", channel, data)
}
