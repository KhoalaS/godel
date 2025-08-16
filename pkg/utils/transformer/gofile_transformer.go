package transformer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/KhoalaS/godel/pkg/auth"
	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

func GofileTransformer(job *types.DownloadJob) error {
	creds, ok := registries.AuthRegistry.Load("gofile")
	if !ok {
		var err error
		creds, err = auth.GofileAuthprovider()
		if err != nil {
			return err
		}
	}

	parsedUrl, err := url.Parse(job.Url)
	if err != nil {
		return err
	}

	id := filepath.Base(parsedUrl.Path)
	dataUrl, err := url.Parse(fmt.Sprintf("https://api.gofile.io/contents/%s", id))
	if err != nil {
		return err
	}

	query := url.Values{}
	query.Add("wt", auth.WT)
	query.Add("page", "1")
	query.Add("pageSize", "1000")
	query.Add("sortField", "name")
	query.Add("sortDirection", "1")

	dataUrl.RawQuery = query.Encode()

	log.Debug().Msgf("Making request to: %s", dataUrl.String())

	request, err := http.NewRequest(http.MethodGet, dataUrl.String(), nil)
	for k, v := range creds.Headers {
		request.Header.Add(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	f, _ := os.Create("gofile_dump.json")
	f.Write(data)

	var fileResponse FileResponse

	err = json.Unmarshal(data, &fileResponse)
	if err != nil {
		return err
	}

	return nil
}

type FileResponse struct {
	Status   string   `json:"status"`
	Data     FileData `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type FileData struct {
	CanAccess          bool                 `json:"canAccess"`
	ID                 string               `json:"id"`
	Type               string               `json:"type"`
	Name               string               `json:"name"`
	CreateTime         int                  `json:"createTime"`
	ModTime            int                  `json:"modTime"`
	Code               string               `json:"code"`
	Public             bool                 `json:"public"`
	TotalDownloadCount int                  `json:"totalDownloadCount"`
	TotalSize          int                  `json:"totalSize"`
	ChildrenCount      int                  `json:"childrenCount"`
	Children           map[string]FileChild `json:"children"`
}

type Metadata struct {
	TotalCount  int  `json:"totalCount"`
	TotalPages  int  `json:"totalPages"`
	Page        int  `json:"page"`
	PageSize    int  `json:"pageSize"`
	HasNextPage bool `json:"hasNextPage"`
}

type FileChild struct {
	CanAccess      bool     `json:"canAccess"`
	ID             string   `json:"id"`
	ParentFolder   string   `json:"parentFolder"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	CreateTime     int      `json:"createTime"`
	ModTime        int      `json:"modTime"`
	Size           int      `json:"size"`
	DownloadCount  int      `json:"downloadCount"`
	Md5            string   `json:"md5"`
	Mimetype       string   `json:"mimetype"`
	Servers        []string `json:"servers"`
	ServerSelected string   `json:"serverSelected"`
	Link           string   `json:"link"`
	Thumbnail      string   `json:"thumbnail"`
}
