package types

type DownloadJob struct {
	Url         string   `json:"url"`
	Filename    string   `json:"filename"`
	Id          string   `json:"id"`
	Password    string   `json:"password"`
	Limit       int      `json:"limit"`
	ConfigId    string   `json:"configId"`
	Transformer []string `json:"transformer"`
	CancelCh    chan struct{}
}
