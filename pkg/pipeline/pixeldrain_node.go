package pipeline

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/KhoalaS/godel/pkg/auth"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

func CreatePixeldrainNode() Node {
	return Node{
		Type: "pixeldrain",
		Run:  PixeldrainNodeFunc,
		Io: map[string]*NodeIO{
			"job": {
				Type:      IOTypeGenerated,
				Id:        "job",
				ValueType: ValueTypeDownloadJob,
				Label:     "Downloader",
				Required:  true,
			},
			"url": {
				Type:      IOTypePassthrough,
				Id:        "url",
				ValueType: ValueTypeString,
				Label:     "Url",
				Required:  true,
			},
		},
		Name:     "Pixeldrain Downloader",
		Status:   StatusPending,
		Category: NodeCategoryDownloader,
	}
}

func isList(pixeldrainUrl string) bool {
	return strings.HasPrefix(pixeldrainUrl, "https://pixeldrain.com/l")
}

func isFile(pixeldrainUrl string) bool {
	return strings.HasPrefix(pixeldrainUrl, "https://pixeldrain.com/u")
}

func PixeldrainNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {

	inputUrl, ok := (node.Io["url"].Value).(string)
	if !ok {
		return errors.New("input url is not a string")
	}

	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Host != "pixeldrain.com" {
		log.Debug().Msg("non pixeldrain url in transformer")
		return errors.New("job url is not a pixeldrain url")
	}

	id := filepath.Base(parsedUrl.Path)

	job := CreateJob(id, inputUrl)

	if job == nil {
		return errors.New("could not create job for pixeldrain url")
	}

	node.Io["job"].Value = job
	return nil
}

func CreateJob(id string, inputUrl string) *types.DownloadJob {
	job := types.NewDownloadJob()

	var fileUrl = ""
	if isFile(inputUrl) {
		fileUrl = fmt.Sprintf("https://pixeldrain.com/api/file/%s?download", id)
	} else if isList(inputUrl) {
		fileUrl = fmt.Sprintf("https://pixeldrain.com/api/list/%s/zip", id)
	} else {
		return nil
	}

	job.Headers["user-agent"] = auth.UserAgent

	job.Url = fileUrl
	return job
}
