package postgres

import (
	"fmt"
	"log"
	"errors"
	"database/sql"
	"github.com/cvenkman/wbL0/internal/config"
)

func Connect(conf config.Config) (*sql.DB, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", conf.DB.Username, conf.DB.Password, conf.DB.Host, conf.DB.Name)
	//?sslmode=disable ???? это что
	// ok = Utils.TryDoIt(s, 10, func() error {
	// 	p.open, ok = sql.Open("postgres", p.connStr)
	// 	return ok
	// })
	fmt.Println(url)
	open, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return open, nil
}

func Select(open *sql.DB, conf config.Config) error {
	q := "SELECT content FROM " + conf.DB.Table + " ;"
	query, err := open.Query(q) // FIXME test from config
	if err != nil {
		return err
	}
	defer func(query *sql.Rows) {
		err := query.Close() // ???
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

func Add(open *sql.DB, modelID string, data []byte, config config.Config) error {
	q := "INSERT INTO " + config.DB.Table + " (id, content) VALUES ($1, $2);"
	_, err := open.Exec(q, modelID, string(data)) // FIXME remove string()
	if err != nil {
		return errors.New("Can't INSERT INTO " + config.DB.Table + ": " + err.Error())
	}
	return nil
}

// func CheckExist(open *sql.DB, modelID string, config config.Config) error {
// 	q := "SELECT id FROM " + config.DB.Table + " WHERE id=" + modelID + ";"
// 	fmt.Println("CheckExist: ", q)
// 	query, err := open.Query(q) // FIXME test from config
// 	if err == nil {
// 		return errors.New("Can't ")
// 	}
// 	query.Close()
// 	return nil
// }