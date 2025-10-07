package pipeline

import "os"

type IFile interface {
	GetDestinationFolder() string
	GetAbsolutePath() (string, error)
	GetFilecontent() ([]byte, error)
	Read(b []byte) (int, error)
	GetFileHandle() (*os.File, error)
}
