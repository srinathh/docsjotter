package main

/*
type NodesServer interface {
	ServeTree(w http.ResponseWriter, r *http.Request)
	ServeNode(w http.ResponseWriter, r *http.Request)
	StartNode(w http.ResponseWriter, r *http.Request)
	EditComment(w http.ResponseWriter, r *http.Request)
}
*/

//FrontEndServer implements the NodeServer Interface for http request handling and defines a back end
//interface for database operations. This loosely couples the client, frontend and backend
type FrontEndServer struct {
}
