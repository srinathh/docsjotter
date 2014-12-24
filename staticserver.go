package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type StaticServer struct {
	root string
}

func NewStaticServer(path string) (StaticServer, error) {
	fi, err := os.Lstat(path)
	if err != nil || !fi.IsDir() {
		return StaticServer{}, fmt.Errorf("NewStaticServer: Bad root path: %s : %s", path, err)
	}
	return StaticServer{path}, nil
}

func (fs StaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var reqpath string
	if r.URL.Path == "/" {
		reqpath = filepath.Join(fs.root, "index.html")
	} else {
		reqpath = filepath.Join(fs.root, r.URL.Path)
	}

	fi, err := os.Lstat(reqpath)
	if err != nil {
		http.Error(w, "StaticServer: Bad Resource Request", http.StatusBadRequest)
		log.Printf("StaticServer: Error finding resource %s : %s", reqpath, err)
		return
	}
	if fi.IsDir() {
		http.Error(w, "StaticServer: Bad Resource Request", http.StatusBadRequest)
		log.Printf("StaticServer: Client requested for directory %s", reqpath)
		return
	}
	http.ServeFile(w, r, reqpath)
}
