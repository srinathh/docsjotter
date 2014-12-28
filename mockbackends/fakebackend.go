package mockbackends

import (
	"fmt"
	"github.com/srinathh/powerdocs/structs"
	"math/rand"

	"github.com/srinathh/powerdocs/utils"
	"sync"
	"time"
)

type FakeBackend struct {
	Locker sync.RWMutex
	Nodes  map[string]structs.Node
}

func (f *FakeBackend) EditComment(id string, c structs.Comment) error {

	f.Locker.Lock()
	defer f.Locker.Unlock()

	val, ok := f.Nodes[id]
	if !ok {
		return fmt.Errorf("FakeBackend.EditComment - Node requested %s not found", id)
	}

	c.ModTime = time.Now().Format(time.Stamp)

	if c.Id == "CREATE" {
		c.Id = utils.DoHash(fmt.Sprintf("%s%d", c.Text, time.Now().UnixNano()))
		val.Comments = append([]structs.Comment{c}, val.Comments...)
		f.Nodes[id] = val
		return nil
	}

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

	rand.Seed(9)
	var f FakeBackend
	f.Locker.Lock()
	f.Nodes = structs.GenNodes()
	f.Locker.Unlock()
	return &f

}
