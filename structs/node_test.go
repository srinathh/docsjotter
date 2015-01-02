package structs

import (
	//"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	n := Node{
		Id:       "0",
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "20 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}
	r, err := http.NewRequest("GET", "/getnode", nil)
	if err != nil {
		t.Fatalf("Error in creating new request :%s", err)
	}
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)

	var n1 Node

	dec := json.NewDecoder(w.Body)
	if err := dec.Decode(&n1); err != nil {
		t.Fatalf("Unable to Decode JSON output :%s \n%s", err, w.Body.String())
	}

	if !reflect.DeepEqual(n1, n) {
		t.Fatalf("Response not matching expectation\n%s", w.Body.String())
	}
}

func TestJSTree(t *testing.T) {
	input := Node{
		Id:       "0",
		Name:     "root",
		Parent:   "#",
		Type:     "directory",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "20 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}

	want := JSTreeNode{
		Id:     "0",
		Parent: "#",
		Text:   "root",
		Type:   "directory",
	}

	got := NewJSTreeNode(input)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("JSTreeNode not as expected:\nWant:\n%v\nGot:\n%v", want, got)
	}

}

/*

func TestReadWriteComments(t *testing.T) {

	want := GenComments(15)

	var buf bytes.Buffer

	if err := WriteComments(want, &buf); err != nil {
		t.Fatalf("Error writing comments:%s", err)
	}

	got, err := ReadComments(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("Error in reading comments : %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Want and got don't match\n%v\n%v", want, got)
	}

}
*/
