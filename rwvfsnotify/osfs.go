package rwvfsnotify

import (
	"fmt"
	"github.com/kr/fs"
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
	pathpkg "path"
	"path/filepath"
	"sourcegraph.com/sourcegraph/rwvfs"
	"strings"
)

type osFS struct {
	rwvfs.WalkableFileSystem

	root         string
	watcher      *fsnotify.Watcher
	done         chan struct{}
	exposeevents chan fsnotify.Event
}

func (osfs *osFS) Events() chan fsnotify.Event {
	return osfs.exposeevents
}

func (osfs *osFS) Errors() chan error {
	return osfs.watcher.Errors
}

//Type FileInfo adds Path attribute to os.FileInfo
type FileInfo struct {
	os.FileInfo
	Path string
}

//func OS returns a virtual FileSystem that uses an fsnotify watcher to watch the
//OS filesystem rooted at root. Since the current implementation of fsnotify does
//not watch folder trees, we walk the path and add add every folder in the tree
//to the watcher. As a convenience since we are scanning the tree, info about all
// the files walked are sent through a chan of rwfsnotify.FileInfo (which adds a
//Path attribute to os.FileInfo) provided by the caller (can be nil)
func OS(root string, infochan chan FileInfo) (FileSystem, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("OS: Error creating FSNotify Watcher:%s", err)
	}

	osfs := &osFS{
		root:               root,
		WalkableFileSystem: rwvfs.Walkable(rwvfs.OS(root)),
		watcher:            watcher,
		done:               make(chan struct{}),
		exposeevents:       make(chan fsnotify.Event),
	}

	go osfs.walk(infochan)
	go osfs.watch()
	return osfs, nil

}

func (osfs *osFS) Done() error {
	close(osfs.done)
	return osfs.watcher.Close()
}

func (osfs *osFS) resolve(path string) string {
	// Clean the path so that it cannot possibly begin with ../.
	// If it did, the result of filepath.Join would be outside the
	// tree rooted at root.  We probably won't ever see a path
	// with .. in it, but be safe anyway.
	path = pathpkg.Clean("/" + path)
	return filepath.Join(osfs.root, path)
}

func (osfs *osFS) watch() {
	for {
		select {
		case <-osfs.done:
			return
		case evt := <-osfs.watcher.Events:
			//strip the root prefix
			path := strings.TrimPrefix(evt.Name, osfs.root)

			if evt.Op == fsnotify.Create {
				info, err := osfs.Stat(path)
				if err != nil {
					log.Printf("Got an fsnotify.Create but unable to stat %s:%s", path, err)
					continue
				}
				if info.IsDir() {
					if err := osfs.watcher.Add(evt.Name); err != nil {
						log.Printf("osfs.watch: got an error adding to watcher fsnotify.Create dir %s:%s", path, err)
					}
				}
			}
			osfs.exposeevents <- fsnotify.Event{path, evt.Op}
		}
	}

}

func (osfs *osFS) walk(infochan chan FileInfo) {
	walker := fs.WalkFS("/", osfs)

	for walker.Step() {
		select {
		case <-osfs.done: //if we receive a done in the middle of loop, just exit
			return
		default:
			if walker.Err() != nil {
				continue
			}

			if walker.Stat().IsDir() {
				if err := osfs.watcher.Add(osfs.resolve(walker.Path())); err != nil {
					log.Printf("osfs.walk: Could not add %s to watcher: %s", walker.Path(), err)
				}
			}

			if infochan != nil {
				infochan <- FileInfo{walker.Stat(), walker.Path()}
			}
		}
	}
}
