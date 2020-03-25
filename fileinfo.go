package memfs

import (
	"os"
	"time"
)

// FileInfo implements os.FileInfo
type FileInfo struct {
	file *File
}

func (s *FileInfo) Name() string {
	return s.file.name
}

func (s *FileInfo) Size() int64 {
	return int64(len(s.file.bytes))
}

func (s *FileInfo) Mode() os.FileMode {
	return os.ModeTemporary
}

func (s *FileInfo) ModTime() time.Time {
	return s.file.modified
}

func (s *FileInfo) IsDir() bool {
	return s.file.bytes == nil
}

func (s *FileInfo) Sys() interface{} {
	return nil
}
