package utils

import (
	"path/filepath"
	"testing"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
)

func TestConfig(t *testing.T) {

	var limitTransformer types.DownloadJobTransformer = func(job *types.DownloadJob) error {
		job.Limit = 1000
		return nil
	}

	registries.TransformerRegistry.Store("limit-transformer-T", limitTransformer)

	config := types.DownloadConfig{
		Id:          "test",
		Name:        "Test-Config",
		DestPath:    "./testfiles",
		Limit:       100,
		Transformer: []string{"limit-transformer-T"},
	}

	job := types.DownloadJob{
		Url: "http://localhost:8080/files/stuff.zip",
	}

	err := ApplyConfig(&job, config)
	if err != nil {
		t.Error(err)
	}

	dest := filepath.Clean(job.Filename)

	if dest != "testfiles/stuff.zip" && dest != "testfiles\\stuff.zip" {
		t.Error(job.Filename)
	}

	if job.Limit != 100 {
		t.Error(job.Limit)
	}
}
