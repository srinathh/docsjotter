package utils

import (
	"encoding/hex"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
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
