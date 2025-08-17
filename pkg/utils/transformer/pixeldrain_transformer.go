package transformer

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/KhoalaS/godel/pkg/auth"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

func PixeldrainTransformer(job *types.DownloadJob) error {
	parsedUrl, err := url.Parse(job.Url)
	if err != nil {
		return err
	}

	if parsedUrl.Host != "pixeldrain.com" {
		log.Debug().Str("url", job.Url).Msg("non pixeldrain url in transformer")
		return errors.New("job url is not a pixeldrain url")
	}

	id := filepath.Base(parsedUrl.Path)

	zipUrl := fmt.Sprintf("https://pixeldrain.com/api/list/%s/zip", id)
	job.Headers["referer"] = job.Url
	job.Headers["user-agent"] = auth.UserAgent

	job.Url = zipUrl

	return nil
}
