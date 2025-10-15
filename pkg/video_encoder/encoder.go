package video_encoder

type VideoEncoder interface {
	Encode(options EncoderOptions) error
}

type EncoderOptions struct {
	OutputFilepath string
	Filename       string
	Container      MediaContainer
	VideoOptions   VideoOptions
	AudioOptions   AudioOptions
}

type VideoOptions struct {
	Bitrate int
	Codec   VideoCodec
}

type AudioOptions struct {
	Bitrate int
	Codec   AudioCodec
}

type VideoCodec = string

const (
	VideoCodecH264  VideoCodec = "libx264"
	VideoCodecH265  VideoCodec = "libx265"
	VideoCodecAV1   VideoCodec = "libaom-av1"
	VideoCodecVP9   VideoCodec = "libvpx-vp9"
	VideoCodecMPEG4 VideoCodec = "mpeg4"
)

type AudioCodec = string

const (
	AudioCodecAAC  AudioCodec = "aac"
	AudioCodecMP3  AudioCodec = "libmp3lame"
	AudioCodecOPUS AudioCodec = "libopus"
	AudioCodecAC3  AudioCodec = "ac3"
	AudioCodecFLAC AudioCodec = "flac"
)

type MediaContainer = string

const (
	MediaContainerMP4  MediaContainer = "mp4"
	MediaContainerMKV  MediaContainer = "mkv"
	MediaContainerWEBM MediaContainer = "webm"
	MediaContainerAVI  MediaContainer = "avi"
	MediaContainerMOV  MediaContainer = "mov"
)
