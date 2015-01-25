//rwvfsnotify is a notifying read/write supported virtual file system. Implemented currently for
//standard OS file systems and may be extended to cover map/memory/zip filesystems
package rwvfsnotify

import (
	"gopkg.in/fsnotify.v1"
	"sourcegraph.com/sourcegraph/rwvfs"
)

//FileSystem encapsulates a read/writable & walkable virtual file system with methods to
//get channels for events and errors from fsnotify which watches the entire directory tree.
//Call the Done() function before terminating to safely close fsnotify watchers. To instantiate
//a new notifying virtual filesystem, call functions like OS
type FileSystem interface {
	rwvfs.WalkableFileSystem
	Events() chan fsnotify.Event
	Errors() chan error
	Done() error
}
