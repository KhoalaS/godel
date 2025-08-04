package utils

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
)

func ApplyConfig(job types.DownloadJob, config types.DownloadConfig) (types.DownloadJob, error) {
	newJob := job

	// apply the transformers
	for _, transformerId := range config.Transformer {
		val, ok := registries.TransformerRegistry.Load(transformerId)
		if !ok {
			continue
		}

		var err error
		newJob, err = val(newJob)
		if err != nil {
			return job, err
		}
	}

	// apply the config values
	if config.Limit > 0 {
		newJob.Limit = config.Limit
	}

	if newJob.Filename == "" {
		inferredName, err := InferFilename(newJob.Url)
		if err != nil {
			return job, err
		}

		newJob.Filename = inferredName
	}

	if config.DestPath != "" {
		newJob.Filename = filepath.Join(config.DestPath, newJob.Filename)
	}

	return newJob, nil
}

func InferFilename(_url string) (string, error) {
	parsedUrl, err := url.Parse(_url)
	if err != nil {
		return "", nil
	}

	spl := strings.Split(parsedUrl.Path, "/")
	return spl[len(spl)-1], nil
}
