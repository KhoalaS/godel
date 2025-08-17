package utils

import (
	"net/url"
	"path"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/google/uuid"
)

func ApplyConfig(job *types.DownloadJob, config types.DownloadConfig) error {
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
		job.DestPath = config.DestPath
	}

	job.DeleteOnCancel = config.DeleteOnCancel

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

	return nil
}

func InferFilename(_url string) (string, error) {
	parsedUrl, err := url.Parse(_url)
	if err != nil {
		return "", err
	}
	name := path.Base(parsedUrl.Path)
	if name == "" {
		name = uuid.NewString()
	}
	return name, nil
}
