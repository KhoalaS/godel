package transformer

import (
	"runtime"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
)

func TestLinuxMountTransformer(t *testing.T) {
	job := types.DownloadJob{
		Filename: "1.txt",
		DestPath: "C:/Users/1/Downloads/",
	}

	LinuxMountTransformer(&job)

	if runtime.GOOS == "linux" {
		if job.DestPath != "/mnt/c/Users/1/Downloads/" {
			t.Fail()
		}
	}
}
