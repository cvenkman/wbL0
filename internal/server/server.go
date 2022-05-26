package server

import (
	"bytes"
	"log"
	// "fmt"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

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
	http.HandleFunc("/order", s.handleOrder)

	err := http.ListenAndServe(conf.Bind_addr, nil)
	if err != nil {
		log.Fatal("serv error", err)
	}
}

// func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
// 	foo, found := s.cache.Get("b2d3d2121d5das63feeb7b")
// 	if found {
// 		io.WriteString(w, foo.(string)) // это что такое
// 	}
// }

func (s *Server) handlfeMain(w http.ResponseWriter, r *http.Request) {
	var data string
	id := "b2d3d2121fdfdgsd5das63feeb7b"
	foo, found := s.cache.Get(id)
	if !found {
		log.Println("handleMain error: " + id + " not found")
		// return
		data = id + " not found"
	} else {
		data = foo.(string)
	}
	// data := "Go Template"
	tmpl, _ := template.New("data").Parse("<h1>{{ .}}</h1>")
	tmpl.Execute(w, data)
}

func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	file, err := template.ParseFiles("internal/server/html/index.html")
	if err != nil {
		log.Fatal(err)
	}
	file.Execute(w, nil)
}

/*TODO соединить index и foo в одну форму*/
 func (s *Server) handleOrder(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	id := r.FormValue("id")
	data, found := s.cache.Get(id)
	if !found {
		io.WriteString(w, "error: id not found")
		log.Println("handleOrder error: " + id + " not found")
		return
		// data = id + " not found"
	}
	if id != "" { // ?? remove id
		tmpl, ok := template.ParseFiles("internal/server/html/foo.html")
		if ok != nil {
			log.Fatal(ok)
		}
		marshal := data.(string)
		var buf bytes.Buffer
		err := json.Indent(&buf, []byte(marshal), "", "\t")
		if err != nil {
			return
		}
		err = tmpl.Execute(w, buf.String())
		if err != nil {
			log.Println("handleOrder error: ", err)
		}
	} else {
		io.WriteString(w, "error: id not found")
	}
	// mutex.Unlock()
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}