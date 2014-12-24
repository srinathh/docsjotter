package main

import (
	"fmt"
)

type FakeBackEnd struct {
	nodes map[string]Node
}

func (f *FakeBackEnd) ListNodeIDs() ([]string, error) {
	ret := make([]string, 0, len(f.nodes))
	for key, _ := range f.nodes {
		ret = append(ret, key)
	}
	return ret, nil
}

func (f *FakeBackEnd) GetNode(id string) (Node, error) {
	if val, ok := f.nodes[id]; !ok {
		return Node{}, fmt.Errorf("FakeBackEnd error in GetNode - couldn't find id: %s", id)
	} else {
		return val, nil
	}
}

func (f *FakeBackEnd) StartNode(id string) error {
	return nil
}

func NewFakeBackend() *FakeBackEnd {
	fakebackend := FakeBackEnd{}
	fakebackend.nodes = make(map[string]Node)
	fakebackend.nodes["0"] = Node{
		Id:       "0",
		Name:     "root",
		Type:     "directory",
		Parent:   "#",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "20 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}

	fakebackend.nodes["1"] = Node{
		Id:       "1",
		Name:     "child",
		Type:     "audio",
		Parent:   "0",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "23 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}

	fakebackend.nodes["2"] = Node{
		Id:       "2",
		Name:     "kitten",
		Type:     "video",
		Parent:   "0",
		ModTime:  "2014-12-19 08:41:00",
		Size:     "23 Mb",
		Comments: []Comment{Comment{"X", "abcdefghijklmnopqrstuvwxyz", "23 Dec 2014"}, Comment{"Y", "1234567890", "22 Dec 2014"}},
	}
	return &fakebackend
}
