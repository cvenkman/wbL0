package main

import (
	"log"
	"os/signal"
	"github.com/nats-io/stan.go"
	// "github.com/nats-io/nats-streaming-server"
	"fmt"
	"sync"
	"os"
)

// func main() {
// 	var clusterID, clientID string
// 	clusterID = "test-cluster"
// 	clientID = "test"
	
// 	sc, err := stan.Connect(clusterID, clientID)
// 	if err != nil {
// 		log.Fatalf("Can't connect publish %s", err.Error())
// 	}
// 	// // Simple Async Subscriber
// 	sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
// 		fmt.Printf("Received a message: %s\n", string(m.Data))
// 	})
	
// 	// for {

// 	// }

// 	// // Unsubscribe
// 	sub.Unsubscribe()
	
// 	// Close connection
// 	sc.Close()
// }

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

func main() {
	var mutex sync.Mutex
	var clusterID, clientID string
	clusterID = "test-cluster"
	clientID = "clientid"
	
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		//log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
		log.Fatalf("Can't connect %s", err.Error())
	}

	mcb := func(msg *stan.Msg) {
		mutex.Lock()
		printMsg(msg, 0)
		mutex.Unlock()
	}

	// _, err = sc.Subscribe(subj, func(msg *stan.Msg) {
	// 	mutex.Lock()
	// 	db.AddNewOrder(&cache, open, msg.Data)
	// 	mutex.Unlock()

	// })

	cleanupDone := make(chan bool)

	// // Simple Async Subscriber
	sub, err := sc.Subscribe("qoo", mcb)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect publish %s", err.Error())
	}
	log.Printf("Listening")
	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			// if durable == "" || unsubscribe {
				sub.Unsubscribe()
			// }
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
