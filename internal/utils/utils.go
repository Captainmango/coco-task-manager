package utils

import (
	"path/filepath"
	"runtime"
)

func BasePath(pathStr string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)

	// go up to project root
	root := filepath.Join(dir, "../..")

	path := filepath.Join(root, pathStr)

	return path
}
