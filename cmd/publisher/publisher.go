package main

import (
	"fmt"
	"log"

	"flag"
	"io/ioutil"

	"github.com/nats-io/stan.go"
)

func main() {
	var clusterID, clientID, channel, dataPath string
	flag.StringVar(&clusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-publisher", "The NATS Streaming client ID to connect with")
	flag.StringVar(&channel, "ch", "test-channel", "The NATS Streaming channel to create")
	flag.StringVar(&dataPath, "data", "model/model.json1", "Json data(model) to publish")
	flag.Parse()

	fmt.Println("Try publish", dataPath, "to channel", channel)

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

	err = sc.Publish(channel, data)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect publish: %s", err.Error())
	}
}
