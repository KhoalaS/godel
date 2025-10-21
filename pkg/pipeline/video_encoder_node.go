package pipeline

import (
	"context"
	"fmt"

	"github.com/KhoalaS/godel/pkg/file"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/KhoalaS/godel/pkg/video_encoder"
)

func NewVideoEncoderNode() Node {
	return Node{
		Type:     "video-encoder",
		Run:      VideoEncoderNodeFunc,
		Name:     "Video Encoder",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"file": {
				Type:      IOTypeConnectedOnly,
				Id:        "file",
				ValueType: ValueTypeFile,
				Label:     "Input File",
			},
			"backend": {
				Type:      IOTypeSelection,
				Id:        "backend",
				ValueType: ValueTypeString,
				Label:     "Backend",
				Options:   []string{"ffmpeg"},
			},
			"videoCodec": {
				Type:      IOTypeSelection,
				Id:        "videoCodec",
				ValueType: ValueTypeString,
				Label:     "Video Codec",
				Options: []video_encoder.VideoCodec{
					video_encoder.VideoCodecCopy,
					video_encoder.VideoCodecH264,
					video_encoder.VideoCodecH265,
					video_encoder.VideoCodecAV1,
					video_encoder.VideoCodecVP9,
					video_encoder.VideoCodecMPEG4,
				},
			},
		},
	}
}

func VideoEncoderNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	file, ok := utils.FromAny[file.IFile](node.Io["file"].Value).Value()
	if !ok {
		return NewInvalidNodeIOError(&node, "file")
	}

	backend, ok := utils.FromAny[string](node.Io["backend"].Value).Value()
	if !ok {
		return NewInvalidNodeIOError(&node, "backend")
	}

	videoCodec, ok := utils.FromAny[video_encoder.VideoCodec](node.Io["videoCodec"].Value).Value()
	if !ok {
		videoCodec = video_encoder.VideoCodecCopy
	}

	var encoder video_encoder.VideoEncoder

	if backend == "ffmpeg" {
		encoder = video_encoder.NewFfmpegEncoder(file)
	}

	if encoder == nil {
		return fmt.Errorf("unkown encoder backend selected: %s", backend)
	}

	// TODO
	encoderOptions := video_encoder.EncoderOptions{
		VideoOptions: video_encoder.VideoOptions{
			Codec: videoCodec,
		},
	}
	err := encoder.Encode(encoderOptions)

	return err
}
