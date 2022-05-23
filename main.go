package main

import (
	"log"
	"os/signal"
	"github.com/nats-io/stan.go"
	"fmt"
	"sync"
	"os"
	"flag"
)

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

func main() {
	var mutex sync.Mutex // ???
	var clusterID, clientID, channel string
	flag.StringVar(&clusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-subscriber", "The NATS Streaming client ID to connect with")
	flag.StringVar(&channel, "ch", "test-channel", "The NATS Streaming channel to create")
	flag.Parse()

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	defer sc.Close()

	msgHandler := func(msg *stan.Msg) { // убп=рать в анонимную функцию внутрь sc.Subscribe
		mutex.Lock() // ???
		printMsg(msg, 0)
		mutex.Unlock()
	}

	
	// // Simple Async Subscriber
	sub, err := sc.Subscribe(channel, msgHandler)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	log.Printf("Listening")
	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	cleanupDone := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}