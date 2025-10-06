package pipeline

type IFile interface {
	GetDestinationFolder() string
	GetAbsolutePath() (string, error)
	GetFilecontent() ([]byte, error)
}
