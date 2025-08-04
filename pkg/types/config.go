package types

type DownloadConfig struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	DestPath      string   `json:"destPath"`
	Transformer   []string `json:"transformer"`
	Limit         int      `json:"limit"`
}
