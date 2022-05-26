package server

import (
	"log"
	// "fmt"
	"net/http"
	"database/sql"
	"io"
	"github.com/cvenkman/wbL0/internal/config"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	conf config.Config
	open *sql.DB
	cache *cache.Cache
}

func New(conf config.Config, open *sql.DB, c *cache.Cache) *Server {
	return &Server {
		conf: conf,
		open: open,
		cache: c,
	}
}

func (s *Server) Serv(conf config.Config, open *sql.DB) {
	// here Start http server listening
	// server.Start()
	//SELECT content FROM delivery WHERE id='b2121d563feb7b';
	log.Println("Server start on port", conf.Bind_addr)

	http.HandleFunc("/", s.handleMain)
	http.HandleFunc("/hello", s.handleHello)

	err := http.ListenAndServe(conf.Bind_addr, nil)
	if err != nil {
		log.Fatal("serv error", err)
	}
}

func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	foo, found := s.cache.Get("b2d3d2121d5das63feeb7b")
	if found {
		io.WriteString(w, foo.(string)) // это что такое
	}
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}