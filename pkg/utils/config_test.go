package utils

import (
	"testing"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
)

func TestConfig(t *testing.T) {

	var limitTransformer types.DownloadJobTransformer = func(job types.DownloadJob) (types.DownloadJob, error) {
		job.Limit = 1000
		return job, nil
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

	newJob, err := ApplyConfig(job, config)
	if err != nil {
		t.Error(err)
	}

	if newJob.Filename != "testfiles\\stuff.zip" && newJob.Filename != "testfiles/stuff.zip" {
		t.Error(newJob.Filename)
	}

	if newJob.Limit != 100 {
		t.Error(newJob.Limit)
	}
}
