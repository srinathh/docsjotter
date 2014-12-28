package mockbackends

import (
	"fmt"
	"github.com/srinathh/powerdocs/structs"
	"reflect"
)

type FailBackend struct {
	Fail    bool
	Comment structs.Comment
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

func (f *FailBackend) EditComment(id string, c structs.Comment) error {
	if f.isFail() != nil {
		return f.isFail()
	}
	if !reflect.DeepEqual(c, f.Comment) {
		return fmt.Errorf("Comment did not match what was expected\nwant\n%v\nGot\n%v\n", f.Comment, c)
	}
	return nil
}
