package pipeline

type Phase int

const (
	PrePhase Phase = iota
	DownloadPhase
	AfterPhase
)
