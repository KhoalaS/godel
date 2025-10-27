package pipeline

import (
	"context"
	"io"
	"net/http"

	"github.com/KhoalaS/godel/pkg/utils"
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
	method, ok := utils.FromAny[string](node.Io["method"].Value).Value()

	if !ok {
		method = http.MethodGet
	}

	url, ok := utils.FromAny[string](node.Io["url"].Value).Value()
	if !ok || url == "" {
		return NewInvalidNodeIOError(&node, "url")
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	node.Io["response"].Value = bodyBytes
	return nil
}
