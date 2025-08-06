package types

import "testing"

func TestJobTransformer(t *testing.T) {
	job := DownloadJob{
		Url:      "a",
		Filename: "b",
		Id:       "0",
	}

	var tr DownloadJobTransformer = func(j *DownloadJob) error {
		j.Filename = "c"
		return nil
	}

	tr(&job)

	if job.Filename != "c" {
		t.Error("old filename not mutated")
	}
}
