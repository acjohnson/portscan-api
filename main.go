package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"github.com/acjohnson/portscan-api/schema"
)

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

func Start(port int) {
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
	return http.StatusOK, "YAY foo"
}

func main() {
	schema.Tables()
	var hosts Hosts
	var scans Scans
	AddResource(hosts, "/hosts")
	AddResource(scans, "/scans")
	Start(4000)
}
