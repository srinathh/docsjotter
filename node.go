package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

//a slice of type JSTreeNode represents the data structure that is sent to JSTree
type JSTreeNode struct {
	Id     string `json:"id"`
	Parent string `json:"parent"`
	Text   string `json:"text"`
	Type   string `json:"type"`
}

func NewJSTreeNode(n Node) JSTreeNode {
	var j JSTreeNode
	j.Id = n.Id
	j.Parent = n.Parent
	j.Text = n.Name
	j.Type = n.Type
	return j
}

//type NodeData represents the data structure served to the client in response to /servenode?id=?
//it is be optionally used by the storage mechanism also (eg. in test backend)
//Comment.Text gets parsed by blackfriday markdown parser before being served
type Node struct {
	Id       string
	Parent   string
	Name     string
	Type     string
	ModTime  string
	Size     string
	Comments []Comment
}

type Comment struct {
	Id      string
	Text    string
	ModTime string
}

func (n Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(n); err != nil {
		http.Error(w, "Could not fetch the node", http.StatusInternalServerError)
		log.Printf("Node.ServeHTTP: Error in encoding node %s to ResponseWriter : %s", n.Name, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, "Could not serve the node", http.StatusInternalServerError)
		log.Printf("Node.ServeHTTP: Error in writing Node %s to ResponseWriter : %s", n.Name, err)
	}

}
