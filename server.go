package main

import (
	"github.com/justinas/alice"
	"log"
	"net/http"
)

//struct FrontEnd stores references to different http.HandlerFunc that will be bound to routes
//by the PowerDocs server. This enables easy swapping of FrontEnd handlers in a piecemeal
//fashion for testing and development. PowerDocs uses github.com/justinas/alice to chain middleware
type FrontEnd struct {
	ServeStatic http.HandlerFunc
	ServeTree   http.HandlerFunc
	ServeNode   http.HandlerFunc
	StartNode   http.HandlerFunc
	EditComment http.HandlerFunc
	MiddleWare  alice.Chain
}

type Config struct {
	IpPort string
}

var config Config
var frontend FrontEnd

func init() {
	frontend.MiddleWare = alice.New(LogHandler)

	if staticserver, err := NewStaticServer("html"); err != nil {
		log.Fatalf("Error configuring static server : %s", err)
	} else {
		frontend.ServeStatic = staticserver.ServeHTTP
	}

	t := NewTestServer()
	backend := NewFakeBackend()
	f := NewFrontEndServer(backend)
	frontend.ServeTree = f.ServeTree
	frontend.ServeNode = f.ServeNode
	frontend.StartNode = f.StartNode
	//	frontend.ServeTree = t.ServeTree
	//	frontend.ServeNode = t.ServeNode
	//	frontend.StartNode = t.StartNode

	frontend.EditComment = t.EditComment

	config.IpPort = "127.0.0.1:8989"
}

func main() {

	//Setting up te resource Handlers
	//http.Handle("/res/", staticserver)
	http.Handle("/", frontend.MiddleWare.ThenFunc(frontend.ServeStatic))
	http.Handle("/servetree", frontend.MiddleWare.ThenFunc(frontend.ServeTree))
	http.Handle("/servenode", frontend.MiddleWare.ThenFunc(frontend.ServeNode))
	http.Handle("/startnode", frontend.MiddleWare.ThenFunc(frontend.StartNode))
	http.Handle("/editcomment", frontend.MiddleWare.ThenFunc(frontend.EditComment))

	//Start the service
	log.Printf("Serving on %s\n", config.IpPort)
	log.Fatal(http.ListenAndServe(config.IpPort, nil))
}

func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s : %s", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
