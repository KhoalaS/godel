package types

import "testing"

func TestJobTransformer(t *testing.T) {
	job := DownloadJob{
		Url:      "a",
		Filename: "b",
		Id:       "0",
	}

	var tr DownloadJobTransformer = func(j DownloadJob) (DownloadJob, error) {
		j.Filename = "c"
		return j, nil
	}

	tr(job)

	if job.Filename != "b" {
		t.Error("mutated filename of old job")
	}
}
