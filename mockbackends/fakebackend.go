package mockbackends

import (
	"fmt"
	"github.com/mailgun/gocql/uuid"
	"github.com/srinathh/docsjotter/structs"
	"github.com/srinathh/docsjotter/utils"

	"log"
	"sync"
	"time"
)

type FakeBackend struct {
	Locker sync.RWMutex
	Nodes  map[string]structs.Node
}

func (f *FakeBackend) EditComment(id, commentid, text string) error {

	f.Locker.Lock()
	defer f.Locker.Unlock()

	log.Printf("Req nodeid: %s   commentid:%s", id, commentid)
	/*
		for k, v := range f.Nodes {
			for _, c := range v.Comments {
				log.Printf("Req nodeid: %s   commentid:%s", k, c.Id)
			}
		}
	*/
	val, ok := f.Nodes[id]
	if !ok {
		return fmt.Errorf("FakeBackend.EditComment - Node requested %s not found", id)
	}
	var c structs.Comment
	c.ModTime = time.Now().Format(utils.MyTimeStamp)
	c.Text = text

	if c.Id == "CREATE" {
		c.Id = uuid.RandomUUID().String()
		val.Comments = append([]structs.Comment{c}, val.Comments...)
		f.Nodes[id] = val
		return nil
	}
	log.Println(val)
	for i, com := range val.Comments {
		if c.Id == com.Id {
			val.Comments[i] = c
			f.Nodes[id] = val
			return nil
		}
	}

	return fmt.Errorf("Could not find the comment to update")
}

func (f *FakeBackend) ListNodes() []structs.JSTreeNode {
	f.Locker.RLock()
	defer f.Locker.RUnlock()
	ret := make([]structs.JSTreeNode, 0, len(f.Nodes))
	for key, value := range f.Nodes {
		ret = append(ret, structs.JSTreeNode{
			Id:     key,
			Type:   value.Type,
			Text:   value.Name,
			Parent: value.Parent,
		})

	}
	return ret
}

func (f *FakeBackend) GetNode(id string) (structs.Node, error) {
	f.Locker.RLock()
	defer f.Locker.RUnlock()
	if val, ok := f.Nodes[id]; !ok {
		return structs.Node{}, fmt.Errorf("FakeBackend error in GetNode - couldn't find id: %s", id)
	} else {
		return val, nil
	}

}

func (f *FakeBackend) StartNode(id string) error {
	return nil
}

func NewFakeBackend() *FakeBackend {

	var f FakeBackend
	f.Locker.Lock()
	f.Nodes = GetTestCases()
	f.Locker.Unlock()
	return &f

}
