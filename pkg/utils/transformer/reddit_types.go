package transformer

type PostsByIdsResponse struct {
	Data PostsInfoByIdsData `json:"data"`
}

type PostsInfoByIdsData struct {
	PostsInfoByIds []PostData `json:"postsInfoByIds"`
}

type PostData struct {
	PostTitle string `json:"postTitle"`
	Media     *Media `json:"media"`
}

type Media struct {
	TypeHint   TypeHint    `json:"typeHint"`
	StillMedia *StillMedia `json:"still"`
	Download   *Download   `json:"download"`
	Streaming  *Streaming  `json:"streaming"`
	Video      *Video      `json:"video"`
	Animated   *Animated   `json:"animated"`
}
type Animated struct {
	TypeName     string       `json:"__typename"`
	Mp4_Source   *MediaSource `json:"mp4_source"`
	Mp4_small    *MediaSource `json:"mp4_small"`
	Mp4_medium   *MediaSource `json:"mp4_medium"`
	Mp4_large    *MediaSource `json:"mp4_large"`
	Mp4_xlarge   *MediaSource `json:"mp4_xlarge"`
	Mp4_xxlarge  *MediaSource `json:"mp4_xxlarge"`
	Mp4_xxxlarge *MediaSource `json:"mp4_xxxlarge"`
	Gif_source   *MediaSource `json:"gif_source"`
	Gif_small    *MediaSource `json:"gif_small"`
	Gif_medium   *MediaSource `json:"gif_medium"`
	Gif_large    *MediaSource `json:"gif_large"`
	Gif_xlarge   *MediaSource `json:"gif_xlarge"`
	Gif_xxlarge  *MediaSource `json:"gif_xxlarge"`
	Gif_xxxlarge *MediaSource `json:"gif_xxxlarge"`
}

type Video struct {
	TypeName   string          `json:"__typename"`
	EmbedHtml  string          `json:"embedHtml"`
	Url        string          `json:"url"`
	Dimensions MediaDimensions `json:"dimensions"`
}

type Streaming struct {
	TypeName         string          `json:"__typename"`
	HlsUrl           string          `json:"hlsUrl"`
	DashUrl          string          `json:"dashUrl"`
	ScrubberMediaUrl string          `json:"scrubberMediaUrl"`
	Dimensions       MediaDimensions `json:"dimensions"`
	Duration         int             `json:"duration"`
	IsGif            bool            `json:"isGif"`
}

type Download struct {
	TypeName string `json:"__typename"`
	Url      string `json:"url"`
}

type StillMedia struct {
	TypeName string      `json:"__typename"`
	Source   MediaSource `json:"source"`
	Small    MediaSource `json:"small"`
	Medium   MediaSource `json:"medium"`
	Large    MediaSource `json:"large"`
	XLarge   MediaSource `json:"xlarge"`
	XXLarge  MediaSource `json:"xxlarge"`
	XXXLarge MediaSource `json:"xxxlarge"`
}

type MediaSource struct {
	TypeName   string          `json:"__typename"`
	Url        string          `json:"url"`
	Dimensions MediaDimensions `json:"dimensions"`
}

type MediaDimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type TypeHint string

const (
	IMAGE    TypeHint = "IMAGE"
	VIDEO    TypeHint = "VIDEO"
	EMBED    TypeHint = "EMBED"
	GIFVIDEO TypeHint = "GIFVIDEO"
)
