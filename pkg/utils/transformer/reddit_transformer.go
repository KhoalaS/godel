package transformer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var fullUrlRegex = regexp.MustCompile(`/r/(.+?)/comments/([^/]+)`)
var shortUrlRegex = regexp.MustCompile(`/r/(.+?)/s/([^/]+)`)

func RedditTransformer(job *types.DownloadJob) error {
	_url, err := url.Parse(job.Url)
	if err != nil {
		return err
	}

	var id string

	match := fullUrlRegex.FindStringSubmatch(_url.Path)
	if len(match) != 3 {
		match = shortUrlRegex.FindStringSubmatch(_url.Path)
		if len(match) != 3 {
			return errors.New("invalid Reddit url")
		}
		id, err = getId(_url)
		if err != nil {
			return err
		}
	} else {
		id = match[2]
	}

	log.Debug().Str("id", id).Msg("extracted id from url")

	creds, ok := registries.AuthRegistry.Load("reddit")
	if !ok {
		log.Debug().Msg("missing credentials, fetching new ones")
		_creds, err := getCredentials()
		if err != nil {
			return err
		}
		log.Debug().Msg("got new credentials")
		creds = *_creds
	} else {
		if creds.Expiry < int(time.Now().Unix()) {
			_creds, err := getCredentials()
			if err != nil {
				return err
			}
			log.Debug().Msg("got new credentials")
			creds.Token = _creds.Token
			creds.Expiry = _creds.Expiry
		}
	}

	post, err := getPostById(id, &creds)
	if err != nil {
		return err
	}

	postInfo := post.Data.PostsInfoByIds[0]

	if postInfo.Media == nil {
		log.Warn().Msg("Post has no media")
		return errors.New("Unimplemented")
	}

	switch postInfo.Media.TypeHint {
	case IMAGE:
		imageUrl := postInfo.Media.StillMedia.Source.Url
		if imageUrl != "" {
			job.Url = imageUrl
		} else {
			return errors.New("source media has no url")
		}
	case EMBED:
		// check if muxed download is available
		if postInfo.Media.Download.Url != "" {
			job.Url = postInfo.Media.Download.Url
			parsedUrl, err := url.Parse(job.Url)
			if err != nil {
				return err
			}
			ext := filepath.Ext(parsedUrl.Path)
			job.Filename = fmt.Sprintf("%s%s", postInfo.PostTitle, ext)
		}
	case GIFVIDEO:
		if postInfo.Media.Animated.Mp4_Source != nil && postInfo.Media.Animated.Mp4_Source.Url != "" {
			job.Url = postInfo.Media.Animated.Mp4_Source.Url
			_, err := url.Parse(job.Url)
			if err != nil {
				return err
			}
			job.Filename = fmt.Sprintf("%s%s", postInfo.PostTitle, ".mp4")
		}
	default:
		log.Warn().Msgf("Post type not implemented yet: %s", postInfo.Media.TypeHint)
		return errors.New("Unimplemented")
	}

	return nil
}

func getCredentials() (*types.Credentials, error) {

	var creds types.Credentials

	body := AuthRequestBody{
		Scopes: []string{"*", "email", "pii"},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/auth/v2/oauth/access-token/loid", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	downRate := float64(20) + rand.Float64()*float64(10)

	request.Header.Add("Authorization", "Basic b2hYcG9xclpZdWIxa2c6")
	request.Header.Add("User-Agent", "Reddit/Version 2025.31.1/Build 2531100/Android 15")
	request.Header.Add("x-reddit-compression", "1")
	request.Header.Add("x-reddit-qos", fmt.Sprintf("down-rate-mbps:%3.f", downRate))
	request.Header.Add("x-reddit-media-codecs", "available-codecs=video/hevc, video/x-vnd.on2.vp9, video/avc")
	request.Header.Add("content-type", "application/json; charset=UTF-8")
	request.Header.Add("client-vendor-id", uuid.NewString())

	response, err := http.DefaultClient.Do(request)
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("got status code != 200")
	}

	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseData AuthResponseBody

	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		return nil, err
	}

	creds = types.Credentials{
		Expiry: responseData.ExpiryTs,
		Token:  responseData.AccessToken,
		Headers: map[string]string{
			"User-Agent":            "Reddit/Version 2025.31.1/Build 2531100/Android 15",
			"x-reddit-compression":  "1",
			"x-reddit-qos":          fmt.Sprintf("down-rate-mbps:%3.f", downRate),
			"x-reddit-media-codecs": "available-codecs=video/hevc, video/x-vnd.on2.vp9, video/avc",
		},
	}

	registries.AuthRegistry.Store("reddit", creds)

	return &creds, nil
}

func getId(_url *url.URL) (string, error) {
	request, err := http.NewRequest(http.MethodGet, _url.String(), nil)
	if err != nil {
		return "", err
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusMovedPermanently {
		return "", errors.New("got wrong status code trying to load short url")
	}

	fullUrl := response.Header.Get("location")
	redirectUrl, err := url.Parse(fullUrl)
	if err != nil {
		return "", err
	}

	match := fullUrlRegex.FindStringSubmatch(redirectUrl.Path)

	if len(match) != 3 {
		return "", errors.New("invalid Reddit url")
	}

	return match[2], nil
}

func getPostById(id string, creds *types.Credentials) (*PostsByIdsResponse, error) {

	requestBody := PostByIdRequest{
		OperationName: "PostsByIds",
		Variables: GQLRequestVariables{
			Ids:                                  []string{fmt.Sprintf("t3_%s", id)},
			IncludeAwards:                        true,
			IncludeEconPromos:                    true,
			IncludeSubredditInPosts:              true,
			IncludePostStats:                     true,
			IncludeDeletedPosts:                  true,
			IncludeCurrentUserAwards:             false,
			IncludeStillMediaAltText:             true,
			IncludeMediaAuth:                     false,
			IncludeExtendedVideoAsset:            false, // ?
			IncludeDevvitData:                    true,
			IncludeCommunityStatus:               true,
			IncludeVideoPlaybackInComments:       false,
			IncludeUnavailablePostReason:         true,
			IncludeModContentDiscussions:         true,
			IncludeNewAdURLField:                 false,
			CommentID:                            "",
			IncludeCommentID:                     false,
			IncludeTranslationDataForDeletedPost: true,
			IncludeSubredditBackgroundColor:      true,
			IncludeAdTransparencyEncodedData:     false,
		},
		Extensions: Extension{
			PersistedQuery: PersistedQuery{
				Version:    1,
				Sha256Hash: "03f9bfc9c4d1d18722fe4f6cbe115dee25377c347e1b4b5d3c97c2a68edc9c10",
			},
		},
	}

	reuqestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://gql-fed.reddit.com/", bytes.NewBuffer(reuqestBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range creds.Headers {
		request.Header.Add(k, v)
	}

	deviceId := uuid.NewString()
	adId := uuid.NewString()

	request.Header.Add("Accept", "application/json")
	request.Header.Add("__temp_suppress_gql_request_latency_seconds", "true")
	request.Header.Add("x-apollo-operation-name", "PostsByIds")
	request.Header.Add("x-apollo-operation-id", "03f9bfc9c4d1d18722fe4f6cbe115dee25377c347e1b4b5d3c97c2a68edc9c10")
	request.Header.Add("x-reddit-http-qos", "LIMITED")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", creds.Token))
	request.Header.Add("client-vendor-id", deviceId)
	request.Header.Add("x-reddit-device-id", deviceId)
	request.Header.Add("x-dev-ad-id", adId)
	request.Header.Add("device-name", "samsung;SM-A5660")
	request.Header.Add("x-reddit-dpr", "4.0")
	request.Header.Add("x-reddit-width", "480")
	request.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	response.Body.Close()

	var obj PostsByIdsResponse
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

type AuthResponseBody struct {
	AccessToken string   `json:"access_token"`
	ExpiryTs    int      `json:"expiry_ts"`
	ExpiresIn   int      `json:"expires_in"`
	Scope       []string `json:"scope"`
	TokenType   string   `json:"token_type"`
}

type AuthRequestBody struct {
	Scopes []string `json:"scopes"`
}

type PostByIdRequest struct {
	OperationName string              `json:"operationName"`
	Variables     GQLRequestVariables `json:"variables"`
	Extensions    Extension           `json:"extensions"`
}

type GQLRequestVariables struct {
	Ids                                  []string `json:"ids"`
	IncludeAwards                        bool     `json:"includeAwards"`
	IncludeEconPromos                    bool     `json:"includeEconPromos"`
	IncludeSubredditInPosts              bool     `json:"includeSubredditInPosts"`
	IncludePostStats                     bool     `json:"includePostStats"`
	IncludeDeletedPosts                  bool     `json:"includeDeletedPosts"`
	IncludeCurrentUserAwards             bool     `json:"includeCurrentUserAwards"`
	IncludeStillMediaAltText             bool     `json:"includeStillMediaAltText"`
	IncludeMediaAuth                     bool     `json:"includeMediaAuth"`
	IncludeExtendedVideoAsset            bool     `json:"includeExtendedVideoAsset"`
	IncludeDevvitData                    bool     `json:"includeDevvitData"`
	IncludeCommunityStatus               bool     `json:"includeCommunityStatus"`
	IncludeVideoPlaybackInComments       bool     `json:"includeVideoPlaybackInComments"`
	IncludeUnavailablePostReason         bool     `json:"includeUnavailablePostReason"`
	IncludeModContentDiscussions         bool     `json:"includeModContentDiscussions"`
	IncludeNewAdURLField                 bool     `json:"includeNewAdUrlField"`
	CommentID                            string   `json:"commentId"`
	IncludeCommentID                     bool     `json:"includeCommentId"`
	IncludeTranslationDataForDeletedPost bool     `json:"includeTranslationDataForDeletedPost"`
	IncludeSubredditBackgroundColor      bool     `json:"includeSubredditBackgroundColor"`
	IncludeAdTransparencyEncodedData     bool     `json:"includeAdTransparencyEncodedData"`
}

type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type Extension struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}
