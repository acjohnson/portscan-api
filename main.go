package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// Test resource
type Test struct {
	ResourceBase
}

// Test GET method
func (t Test) Get(values url.Values) (int, interface{}) {
	return http.StatusOK, "YAY"
}

// Test POST method
func (t Test) Post(values url.Values) (int, interface{}) {
	return http.StatusAccepted, "Post"
}

// Test2 resource
type Test2 struct {
	ResourceBase
}

// Test GET method
func (t Test2) Get(values url.Values) (int, interface{}) {
	return http.StatusOK, "YAY foo"
}

func main() {
	var a Test
	var b Test2
	AddResource(a, "/")
	AddResource(b, "/foo")
	Start(4000)
}

