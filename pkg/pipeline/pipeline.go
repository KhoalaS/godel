package pipeline

type Phase string

const (
	PrePhase      Phase = "pre"
	DownloadPhase Phase = "download"
	AfterPhase    Phase = "after"
)
