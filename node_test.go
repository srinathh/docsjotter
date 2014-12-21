package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/*
var resp = `{
	"Id":"0",
	"Name":"root",
	"Type":"directory",
	"ModTime":"2014-12-19 08:41:00",
	"Size": "20 Mb",
	"Comments":[
		{
			"Id":"X",
			"Text":"abcdefghijklmnopqrstuvwxyz"
		},
		{
			"Id":"X",
			"Text":"1234567890"
		}
	]
}`
*/

func TestNode(t *testing.T) {
	n := Node{
		Id:       "0",
		Name:     "root",
		Type:     "directory",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "20 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz"}, Comment{"Y", "1234567890"}},
	}
	r, err := http.NewRequest("GET", "/getnode", nil)
	if err != nil {
		t.Fatalf("Error in creating new request :%s", err)
	}
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)

	t.Logf("Response:\n%s", w.Body.String())
	var n1 Node

	dec := json.NewDecoder(w.Body)
	if err := dec.Decode(&n1); err != nil {
		t.Fatalf("Unable to Decode JSON output :%s \n%s", err, w.Body.String())
	}

	if !reflect.DeepEqual(n1, n) {
		t.Fatalf("Response not matching expectation\n%s", w.Body.String())
	}
}

func TestTree(t *testing.T) {
	tst := NewTestServer()
	t.Log(tst.Nodes)
}
