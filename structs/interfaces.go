package structs

import (
	"net/http"

	"github.com/justinas/alice"
)

//struct FrontEnd stores references to different http.HandlerFunc that will be bound to routes
//by the PowerDocs server. This enables easy swapping of FrontEnd handlers in a piecemeal
//fashion for testing and development. PowerDocs uses github.com/justinas/alice to chain middleware
type FrontEnd struct {
	ServeStatic http.HandlerFunc
	ServeTree   http.HandlerFunc
	ServeNode   http.HandlerFunc
	StartNode   http.HandlerFunc
	EditComment http.HandlerFunc
	MiddleWare  alice.Chain
}

type BackEnd interface {
	ListNodes() []JSTreeNode
	GetNode(id string) (Node, error)
	StartNode(id string) error
	EditComment(id string, c Comment) error //special ID of CREATE = create new comment
}
