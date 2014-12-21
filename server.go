package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Config struct {
	HTTP     string
	TestMode bool
}

var config Config

func init() {
	flag.StringVar(&config.HTTP, "http", "127.0.0.1:8989", "which host and port to serve on")
	flag.BoolVar(&config.TestMode, "testmode", true, "operate the server in a test configuration with dummy data")
}

type NodesServer interface {
	ServeTree(w http.ResponseWriter, r *http.Request)
	ServeNode(w http.ResponseWriter, r *http.Request)
	StartNode(w http.ResponseWriter, r *http.Request)
}

type ResServer interface {
	ServeRes(w http.ResponseWriter, r *http.Request)
	ServeIndex(w http.ResponseWriter, r *http.Request)
}

func main() {

	//Setup the globals
	flag.Parse()
	rand.Seed(time.Now())

	var nodesserver NodesServer
	var resserver ResServer

	if config.TestMode {
		//t := &TestServer{"html", "testdata"}
		t := NewTestServer()
		nodesserver = t
		resserver = t
	}

	//Setting up te resource Handlers
	http.HandleFunc("/res/", DoLogging(resserver.ServeRes))
	http.HandleFunc("/", DoLogging(resserver.ServeIndex))

	http.HandleFunc("/servetree", DoLogging(nodesserver.ServeTree))
	http.HandleFunc("/servenode", DoLogging(nodesserver.ServeNode))
	//Start the service
	log.Printf("Serving on %s\n", config.HTTP)
	log.Fatal(http.ListenAndServe(config.HTTP, nil))
}

func DoLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.RequestURI())
		h.ServeHTTP(w, r)
	}
}
