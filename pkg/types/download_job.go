package types

import (
	"encoding/json"
	"sync/atomic"
)

type DownloadJob struct {
	Url             string        `json:"url"`
	Filename        string        `json:"filename,omitempty"`
	Id              string        `json:"id"`
	Password        string        `json:"password,omitempty"`
	Limit           int           `json:"limit,omitempty"`
	ConfigId        string        `json:"configId,omitempty"`
	Transformer     []string      `json:"transformer,omitempty"`
	BytesDownloaded int           `json:"bytesDownloaded,omitempty"`
	Size            int           `json:"size,omitempty"`
	DeleteOnCancel  bool          `json:"deleteOnCancel,omitempty"`
	Speed           float64       `json:"speed,omitempty"`
	Eta             float64       `json:"eta,omitempty"`
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
