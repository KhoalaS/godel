package transformer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/KhoalaS/godel/pkg/types"
)

const baseUrl string = "https://api.real-debrid.com/rest/1.0/"

func RealDebridTransformer(job *types.DownloadJob) error {
	apiKey := os.Getenv("RD_KEY")
	if apiKey == "" {
		return fmt.Errorf("no api key")
	}

	_url := baseUrl + "unrestrict/link"

	form := url.Values{}
	form.Add("link", job.Url)
	if job.Password != "" {
		form.Add("password", job.Password)
	}

	req, err := http.NewRequest(http.MethodPost, _url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error in response, code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var unrestrictData UnrestrictResponse

	err = json.Unmarshal(data, &unrestrictData)
	if err != nil {
		return err
	}

	if job.Filename == "" && unrestrictData.Filename != "" {
		job.Filename = unrestrictData.Filename
	}
	job.Url = unrestrictData.Download

	return nil
}

type UnrestrictResponse struct {
	ID         string `json:"id"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mimeType"`
	Filesize   int    `json:"filesize"`
	Link       string `json:"link"`
	Host       string `json:"host"`
	HostIcon   string `json:"host_icon"`
	Chunks     int    `json:"chunks"`
	Crc        int    `json:"crc"`
	Download   string `json:"download"`
	Streamable int    `json:"streamable"`
}
