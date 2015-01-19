package main

import (
	"encoding/json"
	"github.com/srinathh/docsjotter/structs"
	"log"
	"net/http"
)

//FrontEndServer implements the NodeServer Interface for http request handling and defines a back end
//interface for database operations. This loosely couples the client, FrontEndServer and backend
type FrontEndServer struct {
	backend structs.BackEnd
}

func NewFrontEndServer(bk structs.BackEnd) *FrontEndServer {
	var ret FrontEndServer
	ret.backend = bk
	return &ret
}

func (s *FrontEndServer) EditComment(w http.ResponseWriter, r *http.Request) {
	nodeid := r.FormValue("nodeid")
	commentid := r.FormValue("commentid")
	text := r.FormValue("text")

	log.Printf("NodeID: %s, Commentid: %s, text: %s", nodeid, commentid, text)

	if err := s.backend.EditComment(nodeid, commentid, text); err != nil {
		http.Error(w, "Could not edit comment", http.StatusBadRequest)
		log.Printf("Error editing comment: %s : %s", nodeid, err)
		return
	}

}

func (s *FrontEndServer) StartNode(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if err := s.backend.StartNode(id); err != nil {
		http.Error(w, "Could not start requested node", http.StatusBadRequest)
		log.Printf("Error starting requested node : %s : %s", id, err)
		return
	}
}

func (s *FrontEndServer) ServeNode(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	//tk sanitize inputs - allow only hex characterset to pass through

	n, err := s.backend.GetNode(id)
	if err != nil {
		http.Error(w, "Could not find requested node", http.StatusBadRequest)
		log.Printf("Error locating requested node : %s : %s", id, err)
		return
	}
	n.ServeHTTP(w, r)

}

func (t *FrontEndServer) ServeTree(w http.ResponseWriter, r *http.Request) {

	ret := t.backend.ListNodes()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		log.Printf("FrontEndServer.ServeTree: JSON Encoding error:%s", err)
	}
}
