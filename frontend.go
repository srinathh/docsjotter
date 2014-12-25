package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type BackEnd interface {
	ListNodeIDs() ([]string, error)
	GetNode(id string) (Node, error)
	StartNode(id string) error
	EditComment(id string, c Comment) error //special ID of CREATE = create new comment
}

//FrontEndServer implements the NodeServer Interface for http request handling and defines a back end
//interface for database operations. This loosely couples the client, FrontEndServer and backend
type FrontEndServer struct {
	backend BackEnd
}

func NewFrontEndServer(bk BackEnd) *FrontEndServer {
	var ret FrontEndServer
	ret.backend = bk
	return &ret
}

func (s *FrontEndServer) EditComment(w http.ResponseWriter, r *http.Request) {
	nodeid := r.FormValue("nodeid")
	commentid := r.FormValue("commentid")
	text := r.FormValue("text")

	//log.Printf("NodeID: %s, Commentid: %s, text: %s", nodeid, commentid, text)

	//return
	c := Comment{
		Id:      commentid,
		Text:    text,
		ModTime: "",
	}

	if err := s.backend.EditComment(nodeid, c); err != nil {
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
	ids, err := t.backend.ListNodeIDs()
	if err != nil {
		http.Error(w, "Could not fetch list of nodes", http.StatusInternalServerError)
		log.Printf("FrontEndServer.ServeTree: Error in fetching list of node ids:%s", err)
		return
	}

	//making a return slice with capacity of len(ids) and length of 0. There may be errors
	//in fetching some nodes depending on backend and hence avoiding initializing
	ret := make([]JSTreeNode, 0, len(ids))

	for _, id := range ids {
		n, err := t.backend.GetNode(id)
		if err != nil {
			log.Printf("FrontEndServer.ServeTree: Skipping node %s due to error fetching it : %s", id, err)
			continue
		}
		ret = append(ret, NewJSTreeNode(n))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		log.Printf("FrontEndServer.ServeTree: JSON Encoding error:%s", err)
	}
}
