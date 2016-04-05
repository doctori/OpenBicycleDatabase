package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface {
	Get(values url.Values, id int) (int, interface{})
	Post(values url.Values, body io.ReadCloser) (int, interface{})
	Put(values url.Values, body io.ReadCloser) (int, interface{})
	Delete(values url.Values, id int) (int, interface{})
}

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(values url.Values, id int) (int, interface{}) {
	return 405, ""
}

func (PostNotSupported) Post(values url.Values, body io.ReadCloser) (int, interface{}) {
	return 405, ""
}

func (PutNotSupported) Put(values url.Values, body io.ReadCloser) (int, interface{}) {
	return 405, ""
}

func (DeleteNotSupported) Delete(values url.Values, id int) (int, interface{}) {
	return 405, ""
}

type API struct{}

func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}
func (api *API) requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		method := request.Method // Get HTTP Method (string)
		request.ParseForm()      // Populates request.Form
		values := request.Form
		splittedPath := strings.SplitAfter(request.URL.Path, "/")
		log.Print("Received : ")
		log.Println(splittedPath[len(splittedPath)-1])
		id, err := strconv.Atoi(splittedPath[len(splittedPath)-1])
		if err != nil {
			id = 0
		}
		body := request.Body
		fmt.Printf("Received: %s with args : \n\t %+v\n", method, values)
		switch method {
		case "GET":
			code, data = resource.Get(values, id)
		case "POST":
			code, data = resource.Post(values, body)
		case "PUT":
			code, data = resource.Put(values, body)
		case "DELETE":
			code, data = resource.Delete(values, id)
		case "OPTIONS":
			code = 200
			data = nil
		default:
			api.Abort(rw, 405)
		}
		content, err := json.Marshal(data)
		if err != nil {
			api.Abort(rw, 500)
		}
		rw.Header().Set("Content-Type", "text/json; charset=utf-8")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		rw.WriteHeader(code)
		rw.Write(content)
	}
}
func (api *API) AddResource(resource Resource, path string) {
	http.HandleFunc(path, api.requestHandler(resource))
}
func (api *API) Start(inetaddr string, port int) {
	portString := fmt.Sprintf("%s:%d", inetaddr, port)
	log.Fatal(http.ListenAndServe(portString, nil))
}
