package main

import (
	"github.com/justinas/alice"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type NodesServer interface {
	ServeTree(w http.ResponseWriter, r *http.Request)
	ServeNode(w http.ResponseWriter, r *http.Request)
	StartNode(w http.ResponseWriter, r *http.Request)
	EditComment(w http.ResponseWriter, r *http.Request)
}

var nodesserver NodesServer
var staticserver http.Handler
var ipport string
var middleware alice.Chain

func init() {
	var err error
	ipport = "127.0.0.1:8989"

	middleware = alice.New(LogHandler)
	nodesserver = NewTestServer()

	staticserver, err = NewStaticServer("html")
	if err != nil {
		log.Fatalf("Error configuring static server : %s", err)
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	//Setting up te resource Handlers
	//http.Handle("/res/", staticserver)
	http.Handle("/", middleware.Then(staticserver))

	http.Handle("/servetree", middleware.ThenFunc(nodesserver.ServeTree))
	http.Handle("/servenode", middleware.ThenFunc(nodesserver.ServeNode))
	http.Handle("/startnode", middleware.ThenFunc(nodesserver.StartNode))
	http.Handle("/editcomment", middleware.ThenFunc(nodesserver.EditComment))

	//Start the service
	log.Printf("Serving on %s\n", ipport)
	log.Fatal(http.ListenAndServe(ipport, nil))
}

func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s : %s", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
