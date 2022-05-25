package server

import (
	"log"
	"fmt"
	"net/http"
	"database/sql"
	"github.com/cvenkman/wbL0/internal/config"
)

func Serv(conf config.Config, open *sql.DB) {
	// here Start http server listening
	// server.Start()
	//SELECT content FROM delivery WHERE id='b2121d563feb7b';


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "%s", data)

		q := "SELECT content FROM " + conf.DB.Table + " WHERE id='" + "b2121d563feb7b" + "';"
		query, err := open.Query(q) // FIXME test from config
		if err != nil {
			log.Fatal("ser error: ", err)
		}
		for query.Next() {
			var content []byte
			err := query.Scan(&content)
			if err != nil {
				log.Fatal("ser error: ", err)
			}
			fmt.Fprintf(w, "%s", string(content))
		}
	})
	http.ListenAndServe(conf.Bind_addr, nil)
}