package utils

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func DoHash(s string) string {
	hasher := fnv.New64()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func FileSize(root string) int64 {

	var size int64

	log.Printf("Utils.Size: Finding size of :%s", root)

	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {

			log.Printf("Utils.Size: Error in walk: %s for path: %s", err, path)
		} else {
			size = size + fi.Size()
		}
		return nil
	})

	return size
}

func WinStart(pathname string) error {
	log.Printf("Running WinStart on %s\n", pathname)
	cmd := exec.Command("cmd", "/c", "start", "", "/MAX", pathname)
	return cmd.Start()
}

func LinuxStart(pathname string) error {
	log.Printf("Running LinuxStart on %s\n", pathname)
	cmd := exec.Command("xdg-open", pathname)
	return cmd.Start()
}

func Start(pathname string) error {
	switch runtime.GOOS {
	case "linux":
		return LinuxStart(pathname)
	case "windows":
		return WinStart(pathname)
	default:
		return fmt.Errorf("Unsupported operating system %s for starting file", runtime.GOOS)
	}
}
