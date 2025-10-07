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

func (fw *FileWrapper) Read(b []byte) (int, error) {
	return fw.file.Read(b)
}

func (fw *FileWrapper) GetFileHandle() (*os.File, error) {
	fw.file.Close()
	return os.Open(fw.file.Name())
}
