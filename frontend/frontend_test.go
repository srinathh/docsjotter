package frontend

import (
	"encoding/json"
	"github.com/srinathh/powerdocs/mockbackends"
	"github.com/srinathh/powerdocs/structs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestServeNode(t *testing.T) {
	r, err := http.NewRequest("GET", "/servenode?id=15f8c86f8b5e2fac", nil)
	if err != nil {
		t.Fatalf("Error in creating new request :%s", err)
	}
	w := httptest.NewRecorder()

	fakebackend := mockbackends.NewFakeBackend()
	frontend := NewFrontEndServer(fakebackend)

	frontend.ServeNode(w, r)

	dec := json.NewDecoder(w.Body)
	got := structs.Node{}
	if err := dec.Decode(&got); err != nil {
		t.Errorf("Error in decode %s", err)
	}

	want := fakebackend.Nodes["15f8c86f8b5e2fac"]

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

	fakebackend := mockbackends.NewFakeBackend()
	frontend := NewFrontEndServer(fakebackend)

	frontend.ServeTree(w, r)

	dec := json.NewDecoder(w.Body)
	gotslice := []structs.JSTreeNode{}
	if err := dec.Decode(&gotslice); err != nil {
		t.Errorf("Error in decode %s", err)
	}

	for _, got := range gotslice {
		want := structs.NewJSTreeNode(fakebackend.Nodes[got.Id])
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Mismatch in JSTreeNode want and got \n%s\n%s\n", want, got)
		}
	}

}

func TestEditComment(t *testing.T) {

	const text = "### Test Data for comments ###\nChecking if this is `correctly` captured"

	data := url.Values{}
	data.Add("commentid", "CREATE")
	data.Add("text", text)
	data.Add("nodeid", "15f8c86f8b5e2fac")

	r, err := http.NewRequest("POST", "/editcomment", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Error setting up request %s", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

}
