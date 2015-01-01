package structs

import ()

type BackEnd interface {
	ListNodes() []JSTreeNode
	GetNode(id string) (Node, error)
	StartNode(id string) error
	EditComment(id, commentid, text string) error //special ID of CREATE = create new comment
}
