package cache

import (
	"time"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/cvenkman/wbL0/internal/config"
	"database/sql"
)

// Creates a cache and fills it with data from the database
func New(open *sql.DB, conf config.Config) (*cache.Cache, error) {
	c := cache.New(5 * time.Minute, 10 * time.Minute)

	// get data from the database
	q := "SELECT * FROM " + conf.DB.Table + ";"
	query, err := open.Query(q)
	if err != nil {
		return nil, errors.New("Can't select data from " + conf.DB.Table + " table: " + err.Error())
	}

	// add all result to db
	for query.Next() {
		var content, id []byte
		err := query.Scan(&id, &content)
		if err != nil {
			return nil, errors.New("Can't scan query: " + err.Error())
		}
		// add to cache
		c.Set(string(id), string(content), cache.NoExpiration)
	}

	err = query.Close()
	if err != nil {
		return nil, errors.New("Can't close query: " + err.Error())
	}
	return c, nil
}
