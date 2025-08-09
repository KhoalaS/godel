package transformer

type PostsByIdsResponse struct {
	Data PostsInfoByIdsData `json:"data"`
}

type PostsInfoByIdsData struct {
	PostsInfoByIds []PostData `json:"postsInfoByIds"`
}

type PostData struct {
	Media Media `json:"media"`
}

type Media struct {
	TypeHint   TypeHint   `json:"typeHint"`
	StillMedia StillMedia `json:"still"`
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
	IMAGE TypeHint = "IMAGE"
	VIDEO TypeHint = "VIDEO"
)
