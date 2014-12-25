package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type FakeBackend struct {
	Locker sync.RWMutex
	nodes  map[string]Node
}

func (f *FakeBackend) EditComment(id string, c Comment) error {

	f.Locker.Lock()
	defer f.Locker.Unlock()

	val, ok := f.nodes[id]
	if !ok {
		return fmt.Errorf("FakeBackend.EditComment - Node requested %s not found", id)
	}

	c.ModTime = time.Now().Format(time.Stamp)

	if c.Id == "CREATE" {
		c.Id = DoHash(fmt.Sprintf("%s%d", c.Text, time.Now().UnixNano()))
		val.Comments = append([]Comment{c}, val.Comments...)
		f.nodes[id] = val
		return nil
	}

	for i, com := range val.Comments {
		if c.Id == com.Id {
			val.Comments[i] = c
			f.nodes[id] = val
			return nil
		}
	}

	return fmt.Errorf("Could not find the comment to update")
}

func (f *FakeBackend) ListNodeIDs() ([]string, error) {
	f.Locker.RLock()
	defer f.Locker.RUnlock()
	ret := make([]string, 0, len(f.nodes))
	for key, _ := range f.nodes {
		ret = append(ret, key)
	}
	return ret, nil
}

func (f *FakeBackend) GetNode(id string) (Node, error) {
	f.Locker.RLock()
	defer f.Locker.RUnlock()
	if val, ok := f.nodes[id]; !ok {
		return Node{}, fmt.Errorf("FakeBackend error in GetNode - couldn't find id: %s", id)
	} else {
		return val, nil
	}

}

func (f *FakeBackend) StartNode(id string) error {
	return nil
}

func NewFakeBackend() *FakeBackend {

	rand.Seed(9)
	var f FakeBackend
	f.Locker.Lock()
	f.nodes = make(map[string]Node)
	f.genTree()
	f.Locker.Unlock()
	return &f

}

const LOREM string = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

const MAXCOMMENTS = 4

func genComments() []Comment {

	c := make([]Comment, rand.Intn(MAXCOMMENTS))

	for i, _ := range c {
		c[i].Text = LOREM[:rand.Intn(len(LOREM))]
		//adding current nanosecond to prevent comments with same text getting clashing ids
		c[i].Id = DoHash(fmt.Sprintf("%s%d", c[i].Text, time.Now().UnixNano()))
		c[i].ModTime = time.Now().Format(time.Stamp)
	}
	return c

}

func (f *FakeBackend) genTree() {
	f.nodes["0"] = Node{
		Id:       "0",
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["1"] = Node{
		Id:       "1",
		Name:     "Folder 1",
		Type:     "directory",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["2"] = Node{
		Id:       "2",
		Name:     "SubFolder",
		Type:     "directory",
		Parent:   "1",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["3"] = Node{
		Id:       "3",
		Name:     "Audio File in SubFolder",
		Type:     "audio",
		Parent:   "2",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["4"] = Node{
		Id:       "4",
		Name:     "Video File in SubFolder",
		Type:     "video",
		Parent:   "2",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["5"] = Node{
		Id:       "5",
		Name:     "Image File in Folder 1",
		Type:     "image",
		Parent:   "1",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["6"] = Node{
		Id:       "6",
		Name:     "Word File in Folder 1",
		Type:     "msword",
		Parent:   "1",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["7"] = Node{
		Id:       "7",
		Name:     "Folder 2",
		Type:     "directory",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["8"] = Node{
		Id:       "8",
		Name:     "Text File in Folder 2",
		Type:     "text",
		Parent:   "7",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["9"] = Node{
		Id:       "9",
		Name:     "PowerPoint File in Folder 2",
		Type:     "vnd.ms-powerpoint",
		Parent:   "7",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["10"] = Node{
		Id:       "10",
		Name:     "PDF File in Root",
		Type:     "pdf",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}
	f.nodes["11"] = Node{
		Id:       "11",
		Name:     "HTML File in Root",
		Type:     "html",
		Parent:   "0",
		ModTime:  time.Now().Format(time.Stamp),
		Size:     "100 Kb",
		Comments: genComments(),
	}

}

/*
func NewFakeBackend() *FakeBackend {
	FakeBackend := FakeBackend{}
	FakeBackend.nodes = make(map[string]Node)
	FakeBackend.nodes["0"] = Node{
		Id:       "0",
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "20 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}

	FakeBackend.nodes["1"] = Node{
		Id:       "1",
		Name:     "child",
		Type:     "audio",
		Parent:   "0",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "23 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}

	FakeBackend.nodes["2"] = Node{
		Id:       "2",
		Name:     "kitten",
		Type:     "video",
		Parent:   "0",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "23 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}
	return FakeBackend&
}
*/
