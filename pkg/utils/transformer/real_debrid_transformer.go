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

type Endpoint string

const (
	UnrestrictLink Endpoint = "unrestrict/link"
	AddMagnet      Endpoint = "torrents/addMagnet"
)

var endpoints = map[Endpoint]string{
	UnrestrictLink: http.MethodPost,
	AddMagnet:      http.MethodPost,
}

func RealDebridTransformer(job *types.DownloadJob) error {
	params := map[string]string{
		"link": job.Url,
	}
	if job.Password != "" {
		params["password"] = job.Password
	}

	data, err := makeRealDebridRequest[UnrestrictResponse](UnrestrictLink, params)

	if err != nil {
		return err
	}

	if job.Filename == "" && data.Filename != "" {
		job.Filename = data.Filename
	}
	job.Url = data.Download

	return nil
}

func RealDebridMagnetTransformer(job *types.DownloadJob) error {
	params := map[string]string{
		"magnet": job.Url,
	}
	data, err := makeRealDebridRequest[AddMagnetResponse](AddMagnet, params)
	if err != nil {
		return err
	}

	job.Url = data.Uri

	return RealDebridTransformer(job)
}

func makeRealDebridRequest[T any](endpoint Endpoint, params map[string]string) (T, error) {
	var returnValue T

	apiKey := os.Getenv("RD_KEY")
	if apiKey == "" {
		return returnValue, fmt.Errorf("no api key")
	}

	_url := baseUrl + string(endpoint)

	form := url.Values{}

	for k, v := range params {
		form.Add(k, v)
	}

	method := endpoints[endpoint]

	req, err := http.NewRequest(method, _url, strings.NewReader(form.Encode()))
	if err != nil {
		return returnValue, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return returnValue, err
	}

	if response.StatusCode != http.StatusOK {
		return returnValue, fmt.Errorf("error in response, code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return returnValue, err
	}

	err = json.Unmarshal(data, &returnValue)
	if err != nil {
		return returnValue, err
	}

	return returnValue, nil
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

type AddMagnetResponse struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
}
