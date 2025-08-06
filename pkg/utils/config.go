package utils

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
)

func ApplyConfig(job *types.DownloadJob, config types.DownloadConfig) error {
	// apply the transformers
	for _, transformerId := range config.Transformer {
		val, ok := registries.TransformerRegistry.Load(transformerId)
		if !ok {
			continue
		}

		var err error
		err = val(job)
		if err != nil {
			return err
		}
	}

	// apply the config values
	if config.Limit > 0 {
		job.Limit = config.Limit
	}

	if job.Filename == "" {
		inferredName, err := InferFilename(job.Url)
		if err != nil {
			return err
		}

		job.Filename = inferredName
	}

	if config.DestPath != "" {
		job.Filename = filepath.Join(config.DestPath, job.Filename)
	}

	job.DeleteOnCancel = config.DeleteOnCancel

	return nil
}

func InferFilename(_url string) (string, error) {
	parsedUrl, err := url.Parse(_url)
	if err != nil {
		return "", nil
	}

	spl := strings.Split(parsedUrl.Path, "/")
	return spl[len(spl)-1], nil
}
