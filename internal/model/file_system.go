package model

import (
	"io/fs"
)

type FileSystem struct {
	FS      fs.FS
	DirPath string
}

func NewFileSystem(dirPath string, fs fs.FS) FileSystem {
	return FileSystem{FS: fs, DirPath: dirPath}
}
