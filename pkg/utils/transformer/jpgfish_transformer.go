package transformer

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/KhoalaS/godel/pkg/auth"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

var hostRegex = regexp.MustCompile(`jpg[0-9]\.su`)
var albumRegex = regexp.MustCompile(`(?s)class="list-item-image fixed-size.+?class="image-container.+?<img.+?src="(.+?)"`)
var nextPageRegex = regexp.MustCompile(`(?s)data-pagination="next" href="(.+?)"`)

func JpgfishTransformer(job *types.DownloadJob) error {
	// TODO works only for albums
	parsedUrl, err := url.Parse(job.Url)
	if err != nil {
		return err
	}

	if !hostRegex.MatchString(parsedUrl.Host) {
		log.Debug().Str("url", job.Url).Send()
		return errors.New("job url is not a jpg fish url")
	}

	request, err := http.NewRequest(http.MethodGet, job.Url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("referer", job.Url)
	request.Header.Add("user-agent", auth.UserAgent)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	dataString := string(data)
	m := albumRegex.FindAllStringSubmatch(dataString, -1)
	if m == nil {
		return errors.New("could not extract image urls")
	}

	urls := []string{}
	for _, match := range m {
		if len(match) != 2 {
			continue
		}

		image := strings.Replace(match[1], ".md.", ".", 1)
		urls = append(urls, image)
	}

	next := nextPageRegex.FindStringSubmatch(dataString)

	for next != nil {
		if len(next) != 2 {
			break
		}

		nextPage := next[1]
		request, err := http.NewRequest(http.MethodGet, nextPage, nil)
		if err != nil {
			break
		}
		request.Header.Add("referer", job.Url)
		request.Header.Add("user-agent", auth.UserAgent)

		response, err := client.Do(request)
		if err != nil {
			break
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			break
		}

		dataString := string(data)

		m := albumRegex.FindAllStringSubmatch(dataString, -1)
		if m == nil {
			break
		}

		for _, match := range m {
			if len(match) != 2 {
				continue
			}

			image := strings.Replace(match[1], ".md.", ".", 1)
			urls = append(urls, image)
		}

		next = nextPageRegex.FindStringSubmatch(dataString)
	}

	job.IsParent = true
	job.Urls = urls
	job.Headers["user-agent"] = auth.UserAgent

	return nil
}
