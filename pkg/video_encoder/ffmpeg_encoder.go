package video_encoder

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/KhoalaS/godel/pkg/file"
	"github.com/google/uuid"
)

type FfmpegEncoder struct {
	inputFile file.IFile
}

func NewFfmpegEncoder(file file.IFile) *FfmpegEncoder {
	return &FfmpegEncoder{
		inputFile: file,
	}
}

func (e *FfmpegEncoder) Encode(options EncoderOptions) error {
	path, err := e.inputFile.GetAbsolutePath()
	if err != nil {
		return err
	}

	ffmpegCommand := exec.Command("ffmpeg", "-i", path, "-y")

	if options.VideoOptions.Bitrate != 0 {
		ffmpegCommand.Args = append(ffmpegCommand.Args, "-b:v", fmt.Sprintf("%dk", options.VideoOptions.Bitrate))
	}
	if options.AudioOptions.Bitrate != 0 {
		ffmpegCommand.Args = append(ffmpegCommand.Args, "-b:a", fmt.Sprintf("%dk", options.AudioOptions.Bitrate))
	}

	videoCodec := options.VideoOptions.Codec
	if videoCodec == "" {
		videoCodec = "copy"
	}
	ffmpegCommand.Args = append(ffmpegCommand.Args, "-c:v", videoCodec)

	audioCodec := options.AudioOptions.Codec
	if audioCodec == "" {
		audioCodec = "copy"
	}
	ffmpegCommand.Args = append(ffmpegCommand.Args, "-c:a", audioCodec)

	filename := options.Filename
	if filename == "" {
		filename = uuid.NewString()
	}

	filename = fmt.Sprintf("%s.%s", filename, options.Container)
	outputFile := filepath.Join(options.OutputFilepath, filename)
	ffmpegCommand.Args = append(ffmpegCommand.Args, outputFile)

	err = ffmpegCommand.Run()

	return err
}
