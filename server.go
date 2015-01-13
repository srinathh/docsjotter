package main

import (
	"github.com/justinas/alice"
	"github.com/srinathh/powerdocs/fsbackend"
	"github.com/srinathh/powerdocs/mockbackends"
	"github.com/srinathh/powerdocs/structs"
	"log"
	"net/http"
)

//struct FrontEnd stores references to different http.HandlerFunc that will be bound to routes
//by the PowerDocs server. This enables easy swapping of FrontEnd handlers in a piecemeal
//fashion for testing and development. PowerDocs uses github.com/justinas/alice to chain middleware
type AppHandlers struct {
	ServeStatic http.HandlerFunc //HandlerFunc to serve the static assets like javascript and css
	ServeIndex  http.HandlerFunc //HandlerFunc to serve index page
	ServeTree   http.HandlerFunc
	ServeNode   http.HandlerFunc
	StartNode   http.HandlerFunc
	EditComment http.HandlerFunc
	MiddleWare  alice.Chain
}

func main() {

	config := LoadConfig("./docsjotter.toml")

	//setup the middleware
	var apphandlers AppHandlers
	apphandlers.MiddleWare = alice.New(LogHandler)

	//setup the Resource and index Handlers
	staticserver, err := NewStaticServer(config.Resdir)
	if err != nil {
		log.Fatalf("Error configuring static server : %s", err)
	}
	apphandlers.ServeStatic = staticserver.ServeHTTP
	apphandlers.ServeIndex = staticserver.ServeHTTP

	var bk structs.BackEnd

	switch config.Mode {
	case "mocktest":
		bk = mockbackends.NewFakeBackend()
	default:
		bk = fsbackend.NewFSBackend([]string{config.Root})
	}

	f := NewFrontEndServer(bk)

	//we're doing this rather than just use an interface to enable components to be
	//swapped out on the fly
	apphandlers.ServeTree = f.ServeTree
	apphandlers.ServeNode = f.ServeNode
	apphandlers.StartNode = f.StartNode
	apphandlers.EditComment = f.EditComment

	//Setting up te resource Handler	s
	//http.Handle("/res/", staticserver)
	http.Handle("/", apphandlers.MiddleWare.ThenFunc(apphandlers.ServeStatic))
	http.Handle("/servetree", apphandlers.MiddleWare.ThenFunc(apphandlers.ServeTree))
	http.Handle("/servenode", apphandlers.MiddleWare.ThenFunc(apphandlers.ServeNode))
	http.Handle("/startnode", apphandlers.MiddleWare.ThenFunc(apphandlers.StartNode))
	http.Handle("/editcomment", apphandlers.MiddleWare.ThenFunc(apphandlers.EditComment))

	//Start the service
	log.Printf("Serving on %s\n", config.Http)
	log.Fatal(http.ListenAndServe(config.Http, nil))
}

func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s : %s", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
