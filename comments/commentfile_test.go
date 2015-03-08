package comments

import (
	"bytes"
	"github.com/srinathh/hashtag"
	//"reflect"
	"testing"
)

const tomltext = `[ABCDEFGHIJKLMNOP]
  Id = "ABCDEFGHIJKLMNOP"
  BaseName = "test.xlsx"
  ModTime = 2015-02-21T17:49:06Z
  Text = """
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum."""
`

func TestSaveLoad(t *testing.T) {

	var c CommentsFile
	buf := bytes.Buffer{}

	if _, err := c.ReadFrom(bytes.NewReader([]byte(tomltext))); err != nil {
		t.Errorf("Could not load Commentfile: %s", err)
		return
	}

	if _, err := c.WriteTo(&buf); err != nil {
		t.Errorf("Error:%s\nWant:\n%s\nGot:\n%s\n", err, tomltext, buf.Bytes())
		return
	}

	if tomltext != buf.String() {
		t.Errorf("Error:mismatched toml texts\nWant:\n%s\nGot:\n%s\n", tomltext, buf.Bytes())
		return
	}

}

func TestHashTagExtract(t *testing.T) {

	text := "#todo : get fore#cast update #from @alice\nShould we cancel this project? @boss #finance"
	//	wanthashtags := []string{"#todo", "#from", "#finance"}
	//wantattags := []string{"@alice", "@boss"}
	gothashtags := hashtag.ExtractHashtags(text)
	t.Log(gothashtags)
	/*
		if !reflect.DeepEqual(gothashtags, wanthashtags) {
			t.Errorf("Non matching hashtags:\nWant:%v\nGot:%v\n", wanthashtags, gothashtags)
		}
	*/
}
