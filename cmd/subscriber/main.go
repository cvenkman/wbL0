package main

import (
	"log"
	"os/signal"
	"github.com/nats-io/stan.go"
	"fmt"
	"sync"
	"os"
	"flag"
	"github.com/spf13/viper"
	"strings"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq" // <------------ here
	"errors"
)

type Config struct {
	bind_addr string
	db DBconfig
}

type DBconfig struct {
	name string
	table string
	username string
	password string
	host string
}

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

func Connect(config Config) (*sql.DB, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s/%s", config.db.username, config.db.password, config.db.host, config.db.name)
	//?sslmode=disable ???? это что
	fmt.Println(url)
	open, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return open, nil
}

func Select(open *sql.DB, config Config) error {
	q := "SELECT * FROM " + config.db.table + " ;"
	query, err := open.Query(q) // FIXME test from config
	if err != nil {
		return err
	}
	defer func(query *sql.Rows) {
		err := query.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(query)

	for query.Next() {
		var content, id []byte
		err := query.Scan(&id, &content)
		if err != nil {
			return err
		}
		fmt.Println("++ ", string(content))
	}
	return nil
}

func Add(open *sql.DB, data []byte, config Config) error {
	q := "INSERT INTO " + config.db.table + " (id, content) VALUES ($1, $2);"
	_, err := open.Exec(q, 2, string(data)) // FIXME remove string()
	if err != nil {
		return errors.New("Can't INSERT INTO " + config.db.table + ": " + err.Error())
	}
	return nil
}

func readConfig(configPath string) (Config, error) {
	slashIndex := strings.Index(configPath, "/")
	configName := configPath[slashIndex:strings.Index(configPath, ".")]
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath[:slashIndex])
	viper.SetConfigType("toml")

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	config.bind_addr = viper.GetString("bind_addr")
	
	/* get database info from config */
	dbInfo := viper.GetStringMapString("database")
	config.db.name = dbInfo["name"]
	config.db.table = dbInfo["table"]
	config.db.password = dbInfo["password"]
	config.db.host = dbInfo["host"]
	config.db.username = dbInfo["username"]
	return config, nil
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

	config, err := readConfig(configPath)
	if err != nil {
		log.Fatal("Can't read config file: ", err.Error()) // заменитьь log.Fatal на создание своей ошибки и возврт ее из run()
	}

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect to stan %s", err.Error())
	}
	defer sc.Close()

	var str []byte

	open, err := Connect(config)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't connect to %s db %s", config.db.name, err.Error())
	}
	err = Select(open, config)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't SELECT * FROM %s: %s", config.db.table, err.Error())
	}
	/* в этой функции нужно добавить инфу в бд */
	msgHandler := func(msg *stan.Msg) { // убрать в анонимную функцию внутрь sc.Subscribe
		mutex.Lock() // ???
		// printMsg(msg, 0)
		str = msg.Data
		fmt.Println(string(str))
		err = Add(open, str, config)
		if err != nil {
			log.Fatal(err)
		}
		mutex.Unlock()
	}
	
	// // Simple Async Subscriber
	sub, err := sc.Subscribe(channel, msgHandler)
	if err != nil {
		sc.Close()
		log.Fatalf("Can't subscribe to %s channel: %s", channel, err.Error())
	}
	log.Printf("Listening")

	// here Start http server listening
	// server.Start()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(str))
	})
	http.ListenAndServe(config.bind_addr, nil) // FIXME add port from config


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