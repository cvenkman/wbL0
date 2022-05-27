package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cvenkman/wbL0/internal/config"
)

// Connect to postgreSQL
func Connect(conf config.Config) (*sql.DB, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
			conf.DB.Username, conf.DB.Password, conf.DB.Host, conf.DB.Name)

	open, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.New("Can't connect to postgres: " + err.Error())
	}
	fmt.Println("Connect to postgres with url: ", url)
	return open, nil
}

// Add data to postgreSQL
func Add(open *sql.DB, modelID string, data []byte, config config.Config) error {
	q := "INSERT INTO " + config.DB.Table + " (id, content) VALUES ($1, $2);"

	_, err := open.Exec(q, modelID, data)
	if err != nil {
		return errors.New("Can't INSERT INTO " + config.DB.Table + ": " + err.Error())
	}
	return nil
}
