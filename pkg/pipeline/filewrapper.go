package pipeline

import (
	"io"
	"os"
	"path/filepath"
)

type FileWrapper struct {
	file              *os.File
	destinationFolder string
}

func (fw *FileWrapper) GetAbsolutePath() (string, error) {
	return filepath.Abs(fw.file.Name())
}

func (fw *FileWrapper) GetFilecontent() ([]byte, error) {
	return io.ReadAll(fw.file)
}

func (fw *FileWrapper) GetDestinationFolder() string {
	return fw.destinationFolder
}
