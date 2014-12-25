package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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

type FailBackend struct {
	fail    bool
	comment Comment
}

func (f *FailBackend) isFail() error {
	if f.fail {
		return fmt.Errorf("Random Failure in FailBackend")
	}
	return nil
}

func (f *FailBackend) ListNodeIDs() ([]string, error) {
	return []string{}, f.isFail()
}

func (f *FailBackend) GetNode(id string) (Node, error) {
	return Node{}, f.isFail()
}

func (f *FailBackend) StartNode(id string) error {
	return f.isFail()
}

func (f *FailBackend) EditComment(id string, c Comment) error {
	if f.isFail() != nil {
		return f.isFail()
	}
	if !reflect.DeepEqual(c, f.comment) {
		return fmt.Errorf("Comment did not match what was expected\nwant\n%v\nGot\n%v\n", f.comment, c)
	}
	return nil
}

func TestEditComment(t *testing.T) {

	const text = "### Test Data for comments ###\nChecking if this is `correctly` captured"

	data := url.Values{}
	data.Add("commentid", "CREATE")
	data.Add("text", text)
	data.Add("nodeid", "0")

	r, err := http.NewRequest("POST", "/editcomment", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Error setting up request %s", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	failbackend := &FailBackend{false, Comment{Id: "CREATE", ModTime: "", Text: text}}
	frontend := NewFrontEndServer(failbackend)
	frontend.EditComment(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Error in TestEditComment")
	}
}
