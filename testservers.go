package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type TestServer struct {
	HTMLPath     string
	TestDataPath string

	NLock sync.RWMutex
	Nodes map[string]Node
}

func NewTestServer() *TestServer {
	var t TestServer
	t.HTMLPath = "html"
	t.TestDataPath = "testdata"
	t.Nodes = make(map[string]Node)
	t.GenTest(5, 0, 0, "#")
	return &t
}

var testtypes = []string{"audio", "image", "text", "directory"}

const lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin euismod metus sollicitudin turpis cursus mattis eget et eros. In hac habitasse platea dictumst. Donec et orci dolor. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus accumsan felis non eros aliquam sollicitudin. Nulla dolor augue, suscipit nec tortor efficitur, pretium aliquet turpis. Vestibulum egestas consequat massa ac semper. Duis auctor nisi in aliquam viverra. Quisque at malesuada ligula, vel tempor nunc. Aenean placerat dui quis eros congue eleifend. Donec convallis nisi at risus convallis molestie. Duis eu turpis sed nibh iaculis posuere. Duis dapibus velit ultrices vestibulum volutpat. Fusce eget nisl eget diam aliquam blandit non quis libero.\n\nEtiam ultrices ex orci, eget dictum metus blandit a. Suspendisse finibus orci at molestie mollis. Quisque vitae libero at diam accumsan lacinia eu ac neque. Vestibulum maximus justo ut nulla bibendum, vel lobortis tellus porttitor. Ut consectetur facilisis pretium. Maecenas sodales neque vitae pellentesque convallis. Praesent placerat erat orci, sit amet consectetur felis eleifend id. Donec pretium, ex eu porta consequat, nisl felis tempor nibh, nec congue tortor leo et sapien. Ut id urna sit amet orci ultrices sollicitudin eget egestas neque. Aliquam malesuada eu augue ac ultricies. In maximus mattis elementum. Cras quis varius turpis, eu porttitor quam. Nulla id turpis porta, luctus velit vehicula, suscipit lectus. Ut eget ex eu risus accumsan faucibus."
const MAX_NODES int = 10
const MAX_COMMENTS int = 3

//GenTest recursively generates a structure of Nodes with comments to serve on the client test server
func (t *TestServer) GenTest(depth, nodeid, commentid int, parentid string) (int, int) {

	numnodes := rand.Intn(MAX_NODES)
	for j := 0; j < numnodes; j++ {
		var n Node

		n.Id = strconv.FormatInt(int64(nodeid), 10)
		nodeid++

		if depth > 0 {
			n.Type = testtypes[rand.Intn(len(testtypes))]
		} else {
			n.Type = testtypes[rand.Intn(len(testtypes)-1)]
		}

		n.Name = n.Type + " " + n.Id
		n.Parent = parentid
		n.ModTime = time.Now().Format("Jan 2, 2006 at 3:04pm")
		n.Size = n.Id + " kB"

		numcomments := rand.Intn(MAX_COMMENTS)
		n.Comments = make([]Comment, numcomments)
		for i := 0; i < numcomments; i++ {
			c := Comment{"C" + strconv.FormatInt(int64(commentid), 10), lorem[:rand.Intn(len(lorem))], time.Now().Format("Jan 2, 2006 at 3:04pm")}
			commentid++
			n.Comments[i] = c
		}

		t.Nodes[n.Id] = n
		if n.Type == "directory" {
			nodeid, commentid = t.GenTest(depth-1, nodeid, commentid, n.Id)
		}
	}
	return nodeid, commentid
}

func (t *TestServer) ServeTree(w http.ResponseWriter, r *http.Request) {
	t.NLock.RLock()
	ret := make([]JSTreeNode, len(t.Nodes))
	ctr := 0
	for _, v := range t.Nodes {
		ret[ctr] = JSTreeNode{
			Id:     v.Id,
			Parent: v.Parent,
			Text:   v.Name,
			Type:   v.Type,
		}
		ctr++
	}
	t.NLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Printf("Error in TestServer.ServeTree():%s", err)
	}
}

func (t *TestServer) StartNode(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	log.Printf("TestServer.StartNode: Node Start Requested - %s", id)
	t.NLock.RLock()
	defer t.NLock.RUnlock()

	val, ok := t.Nodes[id]
	if !ok {
		http.Error(w, "Couldn't locate requested item", http.StatusBadRequest)
		log.Printf("TestServer.StartNode: Couldn't locate node %s", id)
		return
	}

	success := rand.Intn(5) //each program has 1 in 5 chances of successful start, 0 = failure

	if success == 0 { //failure
		http.Error(w, "Error Starting requested item", http.StatusInternalServerError)
		log.Printf("TestServer.StartNode: Couldn't start node %s", val.Name)
		return
	}
	log.Printf("TestServer.StartNode: Node %s %s successfully started", val.Id, val.Name)
	return //return success

}

func (t *TestServer) ServeNode(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	log.Printf("TestServer.ServeNode: Node Requested - %s", id)

	t.NLock.RLock()
	if val, ok := t.Nodes[id]; ok {
		log.Printf("TestServer.ServeNode: Serving Node %s", val.Name)
		val.ServeHTTP(w, r)
	} else {
		http.Error(w, "Couldn't locate requested item", http.StatusBadRequest)
		log.Printf("TestServer.ServeNode: Couldn't locate node %s", id)
	}
	t.NLock.RUnlock()
}

func (t *TestServer) ServeRes(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/res/", http.FileServer(http.Dir(filepath.Join(t.HTMLPath, "res")))).ServeHTTP(w, r)
}

func (t *TestServer) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(t.HTMLPath, "index.html"))
}
