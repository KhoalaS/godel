package types

type DownloadJob struct {
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Id       int    `json:"id"`
}
