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
	"time"

	"github.com/cvenkman/wbL0/model"
	"github.com/cvenkman/wbL0/internal/config"
	"github.com/cvenkman/wbL0/internal/server"
	"github.com/cvenkman/wbL0/internal/postgres"
	_ "github.com/lib/pq" // <------------ here
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
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

type stanData struct {
	channel		string
	clientID	string
	clusterID	string
}

func parseFlags(configPath *string, stan *stanData) {
	flag.StringVar(configPath, "config", "configs/config.toml", "path to config file")
	flag.StringVar(&stan.clusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&stan.clientID, "id", "stan-subscriber", "The NATS Streaming client ID to connect with")
	flag.StringVar(&stan.channel, "ch", "test-channel", "The NATS Streaming channel to create")
	flag.Parse()
}

func initStan(stanData stanData, msgHandler stan.MsgHandler) (sub stan.Subscription, sc stan.Conn) {
	sc, err := stan.Connect(stanData.clusterID, stanData.clientID)
	if err != nil {
		log.Fatalf("Can't connect to stan %s", err.Error())
	}

	// Simple Async Subscriber
	sub, err = sc.Subscribe(stanData.channel, msgHandler)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't subscribe to %s channel: %s", stanData.channel, err.Error())
	}
	log.Printf("Subscriber listening")

	return
}

func createCache(open *sql.DB, conf config.Config) (c *cache.Cache) {
	c = cache.New(5 * time.Minute, 10 * time.Minute)

	q := "SELECT * FROM " + conf.DB.Table + ";"
	query, err := open.Query(q) // FIXME test from config
	if err != nil {
		log.Fatal("create cache: ", err)
	}
	defer query.Close()

	for query.Next() {
		var content, id []byte
		err := query.Scan(&id, &content)
		if err != nil {
			log.Fatal("create cache: ", err)
		}
		// fmt.Println("++ ", string(content))
		c.Set(string(id), string(content), cache.NoExpiration)
	}
	// c.Set("foo", "bar", cache.DefaultExpiration)
	// c.Set("baz", 42, cache.NoExpiration)

	return
}

func main() {

	var configPath string
	var mutex sync.Mutex // ???
	var stanData stanData
	parseFlags(&configPath, &stanData)

	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal("Can't read config file: ", err.Error()) // заменитьь log.Fatal на создание своей ошибки и возврт ее из run()
	}
	
	open, err := postgres.Connect(config)
	if err != nil {
		log.Fatalf("Can't connect to %s db %s", config.DB.Name, err.Error())
	}

	c := createCache(open, config)

	/* в этой функции нужно добавить инфу в бд */
	msgHandler := func(msg *stan.Msg) { // убрать в анонимную функцию внутрь sc.Subscribe
		mutex.Lock() // ???
		err = addToDB(msg, open, config, c)
		if err != nil {
			log.Println(err)
		}
		mutex.Unlock()
	}
	
	sub, sc := initStan(stanData, msgHandler)

	serverAPI := server.New(config, open, c)
	go serverAPI.Serv(config, open)

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

func addToDB(msg *stan.Msg, open *sql.DB, config config.Config, c *cache.Cache) error {
		// printMsg(msg, 0)
		data := msg.Data
		modelID, err := getModelID(data)
		if err != nil {
			// log.Println(err)
			// mutex.Unlock()
			return err
		}
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
		log.Println("model", modelID, "added")
		// add to cache
		c.Set(modelID, data, cache.NoExpiration)
	return nil
}