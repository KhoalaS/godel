package transformer

import (
	"runtime"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
)

func TestLinuxMountTransformer(t *testing.T) {
	job := types.DownloadJob{
		Filename: "C:/Users/1/Downloads/1.txt",
	}

	LinuxMountTransformer(&job)

	if runtime.GOOS == "linux" {
		if job.Filename != "/mnt/c/Users/1/Downloads/1.txt" {
			t.Fail()
		}
	}
}
