package main

import (
	"github.com/justinas/alice"
	"github.com/srinathh/powerdocs/fsbackend"
	"github.com/srinathh/powerdocs/mockbackends"
	"github.com/srinathh/powerdocs/structs"
	"log"
	"net/http"
)

func main() {

	config := LoadConfig("./docsjotter.toml")

	middleware := alice.New(LogHandler)

	staticserver, err := NewStaticServer(config.Resdir)
	if err != nil {
		log.Fatalf("Error configuring static server : %s", err)
	}

	var bk structs.BackEnd
	switch config.Mode {
	case "mocktest":
		bk = mockbackends.NewFakeBackend()
	default:
		bk = fsbackend.NewFSBackend([]string{config.Root})
	}

	f := NewFrontEndServer(bk)

	http.Handle("/", middleware.Then(staticserver))
	http.Handle("/servetree", middleware.ThenFunc(f.ServeTree))
	http.Handle("/servenode", middleware.ThenFunc(f.ServeNode))
	http.Handle("/startnode", middleware.ThenFunc(f.StartNode))
	http.Handle("/editcomment", middleware.ThenFunc(f.EditComment))

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
