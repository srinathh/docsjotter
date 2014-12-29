package main

import (
	"github.com/BurntSushi/toml"
	"github.com/justinas/alice"
	"github.com/srinathh/powerdocs/frontend"
	"github.com/srinathh/powerdocs/fsbackend"
	//"github.com/srinathh/powerdocs/mockbackends"
	"github.com/srinathh/powerdocs/structs"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Http   string   `toml:"http"`
	Roots  []string `toml:"roots"`
	Resdir string   `toml:"resdir"`
}

var config Config

func SetupConfig() {
	f, err := os.Open("config.toml")
	if err != nil {
		log.Fatalf("Unable to open config.toml : %s", err)
	}

	_, err = toml.DecodeReader(f, &config)
	if err != nil {
		log.Fatalf("Unable to read config.toml:%s", err)
	}
	//tk if some keys are not defined, replace them with default values
}

func main() {
	SetupConfig()
	var frontendserver structs.FrontEnd
	frontendserver.MiddleWare = alice.New(LogHandler)

	if staticserver, err := NewStaticServer(config.Resdir); err != nil {
		log.Fatalf("Error configuring static server : %s", err)
	} else {
		frontendserver.ServeStatic = staticserver.ServeHTTP
	}

	//backend := mockbackends.NewFakeBackend()
	backend := fsbackend.NewFSBackend(config.Roots)
	f := frontend.NewFrontEndServer(backend)

	//we're doing this rather than just use an interface to enable components to be
	//swapped out on the fly
	frontendserver.ServeTree = f.ServeTree
	frontendserver.ServeNode = f.ServeNode
	frontendserver.StartNode = f.StartNode
	frontendserver.EditComment = f.EditComment

	//Setting up te resource Handlers
	//http.Handle("/res/", staticserver)
	http.Handle("/", frontendserver.MiddleWare.ThenFunc(frontendserver.ServeStatic))
	http.Handle("/servetree", frontendserver.MiddleWare.ThenFunc(frontendserver.ServeTree))
	http.Handle("/servenode", frontendserver.MiddleWare.ThenFunc(frontendserver.ServeNode))
	http.Handle("/startnode", frontendserver.MiddleWare.ThenFunc(frontendserver.StartNode))
	http.Handle("/editcomment", frontendserver.MiddleWare.ThenFunc(frontendserver.EditComment))

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
