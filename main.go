package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/acjohnson/portscan-api/database"
	"github.com/acjohnson/portscan-api/logger"
	"github.com/acjohnson/portscan-api/portscanner"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
	"os"
)

var db *sql.DB

type Configuration struct {
	Port     int    `json:"port"`
	LogLevel string `json:"log_level"`
	Debug    bool   `json:"debug"`
	DbType   string `json:"dbtype"`
	DbUri    string `json:"dburi"`
}

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

type Resource interface {
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Delete(values url.Values) (int, interface{})
}

type ResourceBase struct{}

func (ResourceBase) Get(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (ResourceBase) Post(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (ResourceBase) Put(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (ResourceBase) Delete(values url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func requestHandler(resource Resource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data interface{}
		var status int

		r.ParseForm()
		method := r.Method
		values := r.Form

		switch method {
		case MethodGet:
			status, data = resource.Get(values)
		case MethodPost:
			status, data = resource.Post(values)
		case MethodPut:
			status, data = resource.Put(values)
		case MethodDelete:
			status, data = resource.Delete(values)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		w.Write(content)
	}
}

func AddResource(resource Resource, path string) {
	http.HandleFunc(path, requestHandler(resource))
}

func StartServer(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}

// Hosts resource
type Hosts struct {
	ResourceBase
}

// Hosts GET method
func (h Hosts) Get(values url.Values) (int, interface{}) {
	return http.StatusOK, "YAY"
}

// Hosts POST method
func (h Hosts) Post(values url.Values) (int, interface{}) {
	return http.StatusAccepted, "Post"
}

// Scans resource
type Scans struct {
	ResourceBase
}

// Scans GET method
func (s Scans) Get(values url.Values) (int, interface{}) {
	r, err := database.QueryScans(db, values)
	if err != nil {
		log.Fatal(err)
	}
	return http.StatusOK, r
}

// Scans PUT method
func (s Scans) Put(values url.Values) (int, interface{}) {
	port_status, err := portscanner.ScanHost(values.Get("ipv4"))
	r, err := database.UpdateScans(db, values, port_status)
	if err != nil {
		log.Fatal(err)
	}
	return http.StatusAccepted, r
}

func main() {
	var err error

	// Load Config
	conf_file, _ := os.Open("conf.json")
	defer conf_file.Close()
	decoder := json.NewDecoder(conf_file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Set up db connect string
	db, err = sql.Open(configuration.DbType,
		configuration.DbUri)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()

	// Validate db connection
	err = db.Ping()
	if err != nil {
		logger.Println(err)
		panic(err.Error())
	}

	database.Tables(db)
	var hosts Hosts
	var scans Scans
	AddResource(hosts, "/hosts")
	AddResource(scans, "/scans")
	StartServer(configuration.Port)
}
