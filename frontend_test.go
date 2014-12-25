package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServeNode(t *testing.T) {
	r, err := http.NewRequest("GET", "/servenode?id=2", nil)
	if err != nil {
		t.Fatalf("Error in creating new request :%s", err)
	}
	w := httptest.NewRecorder()

	fakebackend := NewFakeBackend()
	frontend := NewFrontEndServer(fakebackend)

	frontend.ServeNode(w, r)

	dec := json.NewDecoder(w.Body)
	got := Node{}
	if err := dec.Decode(&got); err != nil {
		t.Errorf("Error in decode %s", err)
	}

	want := fakebackend.nodes["2"]

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Mismatch in want and got\n%s\n%s", want, got)
	}
}

func TestServeTree(t *testing.T) {
	r, err := http.NewRequest("GET", "/servetree", nil)
	if err != nil {
		t.Fatalf("Error in creating new request :%s", err)
	}
	w := httptest.NewRecorder()

	fakebackend := NewFakeBackend()
	frontend := NewFrontEndServer(fakebackend)

	frontend.ServeTree(w, r)

	dec := json.NewDecoder(w.Body)
	gotslice := []JSTreeNode{}
	if err := dec.Decode(&gotslice); err != nil {
		t.Errorf("Error in decode %s", err)
	}

	for _, got := range gotslice {
		want := NewJSTreeNode(fakebackend.nodes[got.Id])
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Mismatch in JSTreeNode want and got \n%s\n%s\n", want, got)
		}
	}

}
