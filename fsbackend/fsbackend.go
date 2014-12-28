package fsbackend

import (
	"fmt"
	"github.com/srinathh/byten"
	"github.com/srinathh/powerdocs/structs"
	"github.com/srinathh/powerdocs/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FSNode struct {
	structs.JSTreeNode
	path string
}

type FSBackend struct {
	locker  sync.RWMutex
	roots   []string
	fsnodes map[string]FSNode //we will store the path to a directory structure
}

const COMMENTSUFFIX string = ".powerdocs.json"

func NewFSBackend(roots []string) *FSBackend {
	var f FSBackend
	f.roots = roots
	f.fsnodes = make(map[string]FSNode)
	f.populate()
	return &f
}

/*
type BackEnd interface {
	ListNodeIDs() ([]string, error)
	GetNode(id string) (Node, error)
	StartNode(id string) error
	EditComment(id string, c Comment) error //special ID of CREATE = create new comment
}

*/

func (f *FSBackend) StartNode(id string) error {

	f.locker.RLock()
	fsnode, exists := f.fsnodes[id]
	if !exists {
		return fmt.Errorf("FSBackend.StartNode: Error node id %s does not exist", id)
	}
	f.locker.RUnlock()

	if err := utils.Start(fsnode.path); err != nil {
		return fmt.Errorf("FSBackend.Startnode: Error opening node %s:%s", fsnode.path, err)
	}
	return nil
}

func (f *FSBackend) EditComment(id string, c structs.Comment) error {

	f.locker.RLock()
	fsnode, exists := f.fsnodes[id]
	if !exists {
		return fmt.Errorf("FSBackend.EditComment: Error node id %s does not exist", id)
	}
	f.locker.RUnlock()

	isdir := false
	if fsnode.Type == utils.DIRECTORY {
		isdir = true
	}

	comments := fetchcomments(fsnode.path, isdir)

	c.ModTime = time.Now().Format(time.Stamp)

	if c.Id == "CREATE" {
		c.Id = utils.DoHash(fmt.Sprintf("%s%d", c.Text, time.Now().UnixNano()))
		comments = append([]structs.Comment{c}, comments...)
	} else {
		for i, com := range comments {
			if c.Id == com.Id {
				comments[i] = c
				break
			}
		}
	}

	commentpath, exists := getcommentpath(fsnode.path, isdir)

	fil, err := os.Create(commentpath)
	if err != nil {
		return fmt.Errorf("Could not open comment file for writing : %s, %s", commentpath, err)
	}
	defer fil.Close()
	if err := structs.WriteComments(comments, fil); err != nil {
		return fmt.Errorf("Error encoding comments for file : %s, %s", commentpath, err)
	}

	return nil

} //special ID of CREATE = create new comment

func (f *FSBackend) GetNode(id string) (structs.Node, error) {

	f.locker.RLock()
	fsnode, exists := f.fsnodes[id]
	if !exists {
		return structs.Node{}, fmt.Errorf("FSBackend.GetNode: Error node id %s does not exist", id)
	}

	f.locker.RUnlock()

	info, err := os.Lstat(fsnode.path)
	if err != nil {
		return structs.Node{}, fmt.Errorf("Error getting file stats for node id %s: %s", id, err)
	}

	return structs.Node{
		Id:       id,
		Parent:   fsnode.Parent,
		Name:     fsnode.Text,
		Type:     fsnode.Type,
		ModTime:  info.ModTime().Format(time.Stamp),
		Size:     byten.Size(utils.FileSize(fsnode.path)),
		Comments: fetchcomments(fsnode.path, info.IsDir()),
	}, nil
}

func getcommentpath(path string, isdir bool) (string, bool) {
	var commentpath string

	if isdir {
		commentpath = filepath.Join(path, filepath.Base(path)+COMMENTSUFFIX)
	} else {
		commentpath = path + COMMENTSUFFIX
	}

	fi, err := os.Lstat(commentpath)
	if err != nil || fi.IsDir() {
		return commentpath, false
	} else {
		return commentpath, true
	}

}

func fetchcomments(path string, isdir bool) []structs.Comment {

	commentpath, exists := getcommentpath(path, isdir)
	if !exists {
		log.Printf("Comments file doesn't exist for %s or error statting it", path)
		return []structs.Comment{}
	}

	f, err := os.Open(commentpath)
	defer f.Close()
	if err != nil {
		log.Printf("Couldn't open comment file %s:%s", commentpath, err)
		return []structs.Comment{}
	}

	if comments, err := structs.ReadComments(f); err != nil {
		log.Printf("Couldn't parse the comment file %s:%s", commentpath, err)
		return []structs.Comment{}
	} else {

		return comments
	}

}

func (f *FSBackend) ListNodes() []structs.JSTreeNode {
	f.locker.RLock()
	defer f.locker.RUnlock()

	ret := make([]structs.JSTreeNode, 0, len(f.fsnodes))

	for _, fsnode := range f.fsnodes {
		ret = append(ret, fsnode.JSTreeNode)
	}
	return ret
}

func (f *FSBackend) populate() {
	f.locker.Lock()
	defer f.locker.Unlock()

	for _, root := range f.roots {
		err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			//log errors but don't stop walking the path
			if err != nil {
				log.Printf("FSBackend.Walk: Skipping file %s due to error scanning: %s", path, err)
				return nil
			}

			//tk add logic to skip hidden files in dirs

			//skip files which powerdocs stores comments in. They will be checked on the fly
			if strings.HasSuffix(filepath.Base(path), COMMENTSUFFIX) {
				return nil
			}

			var fsnode FSNode
			fsnode.Id = utils.DoHash(path)
			if path == root {
				fsnode.Parent = "#"
			} else {
				fsnode.Parent = utils.DoHash(filepath.Dir(path))
			}
			fsnode.Text = filepath.Base(path)
			fsnode.Type = utils.GetFileType(fi)
			fsnode.path = path
			f.fsnodes[fsnode.Id] = fsnode

			return nil
		})

		if err != nil {
			log.Printf("FSBackend.Walk: Unaccountable error on Walk :%s in root :%s", err, root)
		}
	}
}
