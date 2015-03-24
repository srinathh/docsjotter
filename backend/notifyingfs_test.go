package notifyingvfs

import (
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type Op int

const (
	MKFILE Op = iota
	MKDIR
	RMFILE
	RMDIR
)

var evttexts = map[Op]string{
	MKFILE: "creating file",
	MKDIR:  "creating folder",
	RMFILE: "removing file",
	RMDIR:  "removing folder",
}

func doFileEvent(op Op, name string, t *testing.T, root string) {
	errchk := func(err error) {
		if err != nil {
			t.Fatalf("Error in %s %s: %s", evttexts[op], name, err)
		}
	}

	name = filepath.Join(root, name)

	switch op {
	//case fsnotify.Chmod:
	case MKFILE:
		fil, err := os.Create(name)
		errchk(err)
		fil.Close()

	case RMFILE:
		errchk(os.Remove(name))

	case RMDIR:
		errchk(os.RemoveAll(name))

	case MKDIR:
		errchk(os.Mkdir(name, 0777))

	}

}

func TestNotifyingFs(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	tempdir, err := ioutil.TempDir("", "rwvfsnotify")
	if err != nil {
		t.Fatalf("Error in TestNotifyingFS creating temporary folder :%s", err)
	}
	defer func() {
		if err := os.RemoveAll(tempdir); err != nil {
			t.Logf("Error in removing temp folder %s:%s", tempdir, err)
		}
	}()

	tests := []struct {
		Op   Op
		Path string
		Want fsnotify.Op
	}{
		{
			Op:   MKFILE,
			Path: "/file",
			Want: fsnotify.Create,
		},
		{
			Op:   MKDIR,
			Path: "/folder",
			Want: fsnotify.Create,
		},
		{
			Op:   MKFILE,
			Path: "/folder/file",
			Want: fsnotify.Create,
		},
		{
			Op:   RMFILE,
			Path: "/file",
			Want: fsnotify.Remove,
		},
	}

	osfs, err := OS(tempdir, nil)
	defer osfs.Done()
	time.Sleep(time.Second)

	done := make(chan struct{})
	matches := 0

	go func() {
		var evt fsnotify.Event
		for {
			select {
			case <-done:
				return
			case evt = <-osfs.Events():
				t.Log(evt.String())
				for _, test := range tests {
					if evt.Op == test.Want && evt.Name == test.Path {
						matches++
						t.Log("Found a match")
						break
					}
				}
			}
		}
	}()

	go func() {
		for _, test := range tests {
			t.Logf("%s %s", evttexts[test.Op], test.Path)
			doFileEvent(test.Op, test.Path, t, tempdir)
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second)
		close(done)
	}()
	<-done
	if matches != len(tests) {
		t.Fatalf("Incorrect matches %s", matches)
	}
}
