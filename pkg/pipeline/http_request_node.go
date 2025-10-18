package pipeline

import (
	"context"
	"errors"
	"net/http"
)

func NewHttpRequestNode() Node {
	return Node{
		Type:     "http-request",
		Run:      HttpRequestNodeFunc,
		Name:     "HTTP Request",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"url": {
				Type:      IOTypeInput,
				Id:        "url",
				ValueType: ValueTypeString,
				Label:     "Url",
			},
			"method": {
				Type:      IOTypeSelection,
				Id:        "method",
				ValueType: ValueTypeString,
				Options:   []string{"GET", "POST", "PUT", "HEAD"},
				Label:     "Method",
				Value:     "GET",
			},
			"response": {
				Type:      IOTypeGenerated,
				Id:        "response",
				ValueType: ValueTypeUnknown,
				Label:     "Response",
			},
		},
	}
}

func HttpRequestNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	client := http.DefaultClient
	method := node.Io["method"].Value.(string)

	if method == "" {
		method = http.MethodGet
	}

	url := node.Io["url"].Value.(string)
	if url == "" {
		return errors.New("url was empty")
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	node.Io["response"].Value = res.StatusCode
	return nil
}
