package file

import (
	"io"
	"os"
	"path/filepath"
)

type FileWrapper struct {
	file              *os.File
	destinationFolder string
}

func NewFileWrapper(file *os.File) *FileWrapper {
	return &FileWrapper{
		file: file,
	}
}

func (fw *FileWrapper) GetAbsolutePath() (string, error) {
	return filepath.Abs(fw.file.Name())
}

func (fw *FileWrapper) GetFilecontent() ([]byte, error) {
	return io.ReadAll(fw.file)
}

func (fw *FileWrapper) GetDestinationFolder() (string, error) {
	absPath, err := fw.GetAbsolutePath()
	if err != nil {
		return "", err
	}

	return filepath.Dir(absPath), nil
}

func (fw *FileWrapper) Read(b []byte) (int, error) {
	return fw.file.Read(b)
}

func (fw *FileWrapper) GetFileHandle() (*os.File, error) {
	fw.file.Close()
	return os.Open(fw.file.Name())
}
