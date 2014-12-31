package mockbackends

import (
	"fmt"
	"github.com/srinathh/powerdocs/structs"
)

type FailBackend struct {
	Fail bool
}

func (f *FailBackend) isFail() error {
	if f.Fail {
		return fmt.Errorf("Random Failure in FailBackend")
	}
	return nil
}

func (f *FailBackend) ListNodes() []structs.JSTreeNode {
	return []structs.JSTreeNode{}
}

func (f *FailBackend) GetNode(id string) (structs.Node, error) {
	return structs.Node{}, f.isFail()
}

func (f *FailBackend) StartNode(id string) error {
	return f.isFail()
}

func (f *FailBackend) EditComment(id, commentid, text string) error {
	if f.isFail() != nil {
		return f.isFail()
	}

	return nil
}
