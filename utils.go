package main

import (
	"encoding/hex"
	"hash/fnv"
)

func DoHash(s string) string {
	hasher := fnv.New64()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
