package main

import (
	"github.com/justinas/alice"
	"github.com/srinathh/powerdocs/frontend"
	"github.com/srinathh/powerdocs/fsbackend"
	//"github.com/srinathh/powerdocs/mockbackends"
	"github.com/srinathh/powerdocs/structs"
	"log"
	"net/http"
)

type Config struct {
	IpPort string
}

var config Config
var frontendserver structs.FrontEnd

func init() {
	frontendserver.MiddleWare = alice.New(LogHandler)

	if staticserver, err := NewStaticServer("html"); err != nil {
		log.Fatalf("Error configuring static server : %s", err)
	} else {
		frontendserver.ServeStatic = staticserver.ServeHTTP
	}

	//backend := mockbackends.NewFakeBackend()
	backend := fsbackend.NewFSBackend([]string{"testdata"})
	f := frontend.NewFrontEndServer(backend)
	frontendserver.ServeTree = f.ServeTree
	frontendserver.ServeNode = f.ServeNode
	frontendserver.StartNode = f.StartNode
	frontendserver.EditComment = f.EditComment

	config.IpPort = "127.0.0.1:8989"
}

func main() {

	//Setting up te resource Handlers
	//http.Handle("/res/", staticserver)
	http.Handle("/", frontendserver.MiddleWare.ThenFunc(frontendserver.ServeStatic))
	http.Handle("/servetree", frontendserver.MiddleWare.ThenFunc(frontendserver.ServeTree))
	http.Handle("/servenode", frontendserver.MiddleWare.ThenFunc(frontendserver.ServeNode))
	http.Handle("/startnode", frontendserver.MiddleWare.ThenFunc(frontendserver.StartNode))
	http.Handle("/editcomment", frontendserver.MiddleWare.ThenFunc(frontendserver.EditComment))

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
