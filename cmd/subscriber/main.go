package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"encoding/json"
	"database/sql"

	"github.com/cvenkman/wbL0/model"
	"github.com/cvenkman/wbL0/internal/config"
	"github.com/cvenkman/wbL0/internal/server"
	"github.com/cvenkman/wbL0/internal/postgres"
	_ "github.com/lib/pq" // <------------ here
	"github.com/nats-io/stan.go"
)

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

func getModelID(data []byte) (string, error) {
	var model model.Delivery
	err := json.Unmarshal(data, &model)
	if err != nil || model.OrderUID == "" {
		return "", errors.New("Unmarshal error: ")
	}
	return model.OrderUID, nil
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.toml", "path to config file")
	
	var mutex sync.Mutex // ???

	var clusterID, clientID, channel string
	flag.StringVar(&clusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-subscriber", "The NATS Streaming client ID to connect with")
	flag.StringVar(&channel, "ch", "test-channel", "The NATS Streaming channel to create")
	flag.Parse()

	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal("Can't read config file: ", err.Error()) // заменитьь log.Fatal на создание своей ошибки и возврт ее из run()
	}

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect to stan %s", err.Error())
	}
	defer sc.Close()

	// var data []byte

	open, err := postgres.Connect(config)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect to %s db %s", config.DB.Name, err.Error())
	}
	/* в этой функции нужно добавить инфу в бд */
	msgHandler := func(msg *stan.Msg) { // убрать в анонимную функцию внутрь sc.Subscribe
		mutex.Lock() // ???
		err = addToDB(msg, open, config)
		if err != nil {
			log.Println(err)
		}
		mutex.Unlock()
	}
	
	// Simple Async Subscriber
	sub, err := sc.Subscribe(channel, msgHandler)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't subscribe to %s channel: %s", channel, err.Error())
	}
	log.Printf("Listening")

	go server.Serv(config, open)

	cleanupDone := make(chan bool)
	cleanupAfterSIGINT(cleanupDone, sub, sc)
	<-cleanupDone
}

// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
func cleanupAfterSIGINT(cleanupDone chan bool, sub stan.Subscription, sc stan.Conn) {
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
}

func addToDB(msg *stan.Msg, open *sql.DB, config config.Config) error {
		// printMsg(msg, 0)
		data := msg.Data
		modelID, err := getModelID(data)
		if err != nil {
			// log.Println(err)
			// mutex.Unlock()
			return err
		}
		fmt.Println(modelID)
		// err = CheckExist(open, modelID, config)
		// if err != nil {
		// 	log.Println(err)
		// 	mutex.Unlock()
		// 	return
		// }
		err = postgres.Add(open, modelID, data, config)
		if err != nil {
			// log.Println(err)
			return err
		}
	return nil
}