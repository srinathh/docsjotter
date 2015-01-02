package main

import (
	"fmt"
	"github.com/cespare/flagconf"
	"github.com/justinas/alice"
	"github.com/srinathh/powerdocs/frontend"
	"github.com/srinathh/powerdocs/fsbackend"
	"github.com/srinathh/powerdocs/mockbackends"
	"github.com/srinathh/powerdocs/structs"
	"log"
	"net/http"
)

type Config struct {
	Http   string `desc:"the IP address and port for the DocsJotter server. Defaults to 127.0.0.1:8989"`
	Root   string `desc:"which directory to serve"`
	Resdir string `desc:"path to the DocsJotter resource files. Defaults to html"`
	Mode   string `desc:"can take values production, fstest or mocktest. In production mode, Root must be speciried. In fstest, root defaults to testdata. MockTest uses a mock backend"`
}

var config Config

func ParseConfig() error {
	config = Config{
		Http:   "127.0.0.1:8989",
		Root:   "",
		Resdir: "html",
		Mode:   "fstest",
	}

	if err := flagconf.Parse("docsjotter.toml", &config); err != nil {
		return fmt.Errorf("Could not find docsjotter.toml: %s", err)
	}

	switch config.Mode {
	case "fstest":
		config.Root = "testdata"
	}
	return nil
}

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
	if err := ParseConfig(); err != nil {
		log.Fatalf("Error reading configuration %s:", err)
	}

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

	f := frontend.NewFrontEndServer(bk)

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
