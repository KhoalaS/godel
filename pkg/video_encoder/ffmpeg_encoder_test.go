package video_encoder

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/KhoalaS/godel/pkg/file"
)

func TestFfmpegEncoder(t *testing.T) {
	encoderOptions := EncoderOptions{
		OutputFilepath: "./",
		Filename:       "outfile",
		Container:      MediaContainerMP4,
	}

	fp, _ := filepath.Abs("../../testfiles/Main.mp4")
	testFile, err := os.Open(fp)

	if err != nil {
		t.Error(err)
	}

	fileWrapper := file.NewFileWrapper(testFile)

	ffmpegEncoder := NewFfmpegEncoder(fileWrapper)
	err = ffmpegEncoder.Encode(encoderOptions)
	if err != nil {
		t.Error(err)
	}

	testFile.Close()
	os.Remove(fmt.Sprintf("./outfile.%s", encoderOptions.Container))
}
