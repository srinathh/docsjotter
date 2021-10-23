//package backend implements file system (or pseudo file system) access
//operations. adds file system notifications to SourceGraph's
//read/write virtual file system package. Currently implemented for OsFs.
package backend

import (
	"path/filepath"

	"github.com/rjeczalik/notify"
	"sourcegraph.com/sourcegraph/rwvfs"
)

//type FileSystem defines a read/write file system that also allows
//event monitoring in line with the rjeczalik/notify framework
type FileSystem interface {
	rwvfs.WalkableFileSystem
	Stop(c chan<- notify.EventInfo)
	Watch(path string, c chan<- notify.EventInfo, events ...notify.Event) error
}

type OsFS struct {
	rwvfs.WalkableFileSystem
	root string
}

//func OS returns a read/write and notify supporting walkable virtual FileSystem
//rooted at the parameter root on the underlying file system provided by OS
func OS(root string) FileSystem {
	return &OsFS{
		root:               root,
		WalkableFileSystem: rwvfs.Walkable(rwvfs.OS(root)),
	}
}

// resolve is from golang.org/x/tools/godoc/vfs.
func (fs *OsFS) resolve(path string) string {
	// Clean the path so that it cannot possibly begin with ../.
	// If it did, the result of filepath.Join would be outside the
	// tree rooted at root.  We probably won't ever see a path
	// with .. in it, but be safe anyway.
	path = pathpkg.Clean("/" + path)
	return filepath.Join(string(fs.root), path)
}

func (osfs *OsFs) Watch(path string, c chan<- notify.EventInfo, events ...notify.Event) error {
	return notify.Watch(osfs.resolve(path), c, events)
}

func (osfs *OsFs) Stop(c chan<- notify.EventInfo) {
	return notify.Stop(c)

}

/*
FileSystem functions
    Open(name string) (ReadSeekCloser, error)
    Lstat(path string) (os.FileInfo, error)
    Stat(path string) (os.FileInfo, error)
    ReadDir(path string) ([]os.FileInfo, error)
    String() string
    Create(path string) (io.WriteCloser, error)
    Mkdir(name string) error
    Remove(name string) error

Create, Remove, Write and Rename are the only event values guaranteed to be present on all platforms.


type MapFs struct {
	fs     rwvfs.WalkableFileSystem
	events map[string][]watch
}

type watch struct {
	event notify.Event
	c     chan<- notify.EventInfo
}

//func Map returns a read/write and notify supporting walkable virtual FileSystem
//based on the provided map[string]string parameter. Keys should be forward slash
//separated paths and values should be the content of the files
func Map(m map[string]string) FileSystem {
	return MapFs{
		fs: rwvfs.WalkableFileSystem(rwvfs.Map(m)),
	}
}

func (mapfs *MapFs) Watch(path string, c chan<- notify.EventInfo, events ...notify.Event) error {
	for k, v := range mapfs.events{
		if filepath.HasPrefix(p, prefix)
	}
}

func (osfs *OsFs) Stop(c chan<- notify.EventInfo) {
	return notify.Stop(c)
}
*/
