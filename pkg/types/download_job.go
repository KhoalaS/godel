package types

import (
	"encoding/json"
	"sync/atomic"
)

type DownloadJob struct {
	Url             string        `json:"url"`
	Filename        string        `json:"filename"`
	Id              string        `json:"id"`
	Password        string        `json:"password"`
	Limit           int           `json:"limit"`
	ConfigId        string        `json:"configId"`
	Transformer     []string      `json:"transformer"`
	BytesDownloaded int           `json:"bytesDownloaded"`
	Size            int           `json:"size"`
	DeleteOnCancel  bool          `json:"deleteOnCancel"`
	Status          atomic.Value  `json:"-"`
	CancelCh        chan struct{} `json:"-"`
	PauseCh         chan struct{} `json:"-"`
}

func (j *DownloadJob) MarshalJSON() ([]byte, error) {
	type Alias DownloadJob // prevent recursion
	return json.Marshal(&struct {
		Status string `json:"status"`
		*Alias
	}{
		Status: func() string {
			val := j.Status.Load()
			if val == nil {
				return ""
			}
			return string(val.(DownloadState))
		}(),
		Alias: (*Alias)(j),
	})
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
