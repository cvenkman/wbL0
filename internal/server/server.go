package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/cvenkman/wbL0/internal/config"
	"github.com/patrickmn/go-cache"
)

// Struct for server
type Server struct {
	conf config.Config
	open *sql.DB
	cache *cache.Cache
}

type ViewData struct {
	Model string
	Lenght string
}

// Create server struct
func New(conf config.Config, open *sql.DB, c *cache.Cache) *Server {
	return &Server {
		conf: conf,
		open: open,
		cache: c,
	}
}

// Start http server listening
func (s *Server) Start() {
	log.Println("Server start on port", s.conf.Bind_addr)

	http.HandleFunc("/", s.handleMain)
	http.HandleFunc("/order", s.handleOrder)

	err := http.ListenAndServe(s.conf.Bind_addr, nil)
	if err != nil {
		log.Fatal("serv error", err)
	}
}

// Handle for main page
func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	lenght := strconv.Itoa(s.cache.ItemCount())
	file, err := template.ParseFiles("internal/server/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	file.Execute(w, lenght)
}

// Handle with form for order input
 func (s *Server) handleOrder(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	data, found := s.cache.Get(id)
	if !found {
		errorNotFound := "Can't find record with id " + id
		log.Println(errorNotFound)

		templateErr, err := template.ParseFiles("internal/server/templates/modelNotFound.html")
		err = templateErr.Execute(w, errorNotFound)
		if err != nil {
			log.Println("handleOrder error: ", err)
		}
		return
	}

	template, err := template.ParseFiles("internal/server/templates/model.html")
	if err != nil {
		log.Fatal(err)
	}

	// indent json file
	marshal := data.(string)
	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(marshal), "", "\t")
	if err != nil {
		return
	}

	err = template.Execute(w, buf.String())
	if err != nil {
		log.Println("Can't execute template with json: ", err)
	}
}
