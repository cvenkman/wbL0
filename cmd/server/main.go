package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"

	dbCache "github.com/cvenkman/wbL0/internal/cache"
	"github.com/cvenkman/wbL0/internal/config"
	"github.com/cvenkman/wbL0/internal/postgres"
	"github.com/cvenkman/wbL0/internal/server"
	"github.com/cvenkman/wbL0/internal/stan"
	"github.com/cvenkman/wbL0/internal/utils"
	"github.com/cvenkman/wbL0/model"
	_ "github.com/lib/pq"
	stanAPI "github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
)

func parseFlags(configPath *string, st *stan.Stan) {
	flag.StringVar(configPath, "config", "configs/config.toml", "path to config file")
	flag.StringVar(&st.ClusterID, "cl", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&st.ClientID, "id", "stan-subscriber", "The NATS Streaming client ID to connect with")
	flag.StringVar(&st.Channel, "ch", "test-channel", "The NATS Streaming channel to create")
	flag.Parse()
}

// read config, connect to postgres, create cache and start server
// in the end waiting for a SIGINT and cleanup
func main() {
	var configPath string
	var st stan.Stan
	parseFlags(&configPath, &st)

	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	
	open, err := postgres.Connect(config)
	if err != nil {
		log.Fatal(err)
	}
	
	cache, err := dbCache.New(open, config)
	if err != nil {
		log.Fatal(err)
	}
	
	sub, sc := st.InitStan(func(msg *stanAPI.Msg) {
		err = saveData(msg.Data, open, config, cache)
		msg.Ack()
		if err != nil {
			log.Println(err)
		}
	})

	serverAPI := server.New(config, open, cache)
	serverAPI.Start()

	cleanupDone := make(chan bool)
	utils.CleanupAfterSIGINT(cleanupDone, sub, sc)
	<-cleanupDone
}

// add data from json to db and cache
func saveData(data []byte, open *sql.DB, config config.Config, c *cache.Cache) error {
	modelID, err := model.GetID(data)
	if err != nil {
		return err
	}
	if modelID == "" {
		return errors.New("Error: empty ID")
	}

	err = postgres.Add(open, modelID, data, config)
	if err != nil {
		return err
	}
	c.Set(modelID, string(data), cache.NoExpiration)

	log.Println("model", modelID, "added")
	return nil
}
