package file

import "os"

type IFile interface {
	GetDestinationFolder() (string, error)
	GetAbsolutePath() (string, error)
	GetFilecontent() ([]byte, error)
	Read(b []byte) (int, error)
	GetFileHandle() (*os.File, error)
}
