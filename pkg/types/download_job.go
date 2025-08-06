package types

import "sync/atomic"

type DownloadJob struct {
	Url             string       `json:"url"`
	Filename        string       `json:"filename"`
	Id              string       `json:"id"`
	Password        string       `json:"password"`
	Limit           int          `json:"limit"`
	ConfigId        string       `json:"configId"`
	Transformer     []string     `json:"transformer"`
	BytesDownloaded int          `json:"bytesDownloaded"`
	Status          atomic.Value `json:"status"`
	CancelCh        chan struct{}
	PauseCh         chan struct{}
}

func NewDownloadJob() *DownloadJob {
	job := DownloadJob{CancelCh: make(chan struct{}), PauseCh: make(chan struct{}), Transformer: []string{}}
	job.Status.Store(IDLE)
	return &job
}

type DownloadState string

const (
	IDLE        DownloadState = "idle"
	PAUSED      DownloadState = "paused"
	CANCELED    DownloadState = "canceled"
	DOWNLOADING DownloadState = "downloading"
	DONE        DownloadState = "done"
	ERROR       DownloadState = "error"
)
