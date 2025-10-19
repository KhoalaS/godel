package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/services"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/rs/zerolog/log"
)

func NewAddTorrentNode() Node {
	return Node{
		Type:     "add-torrent",
		Name:     "Add Torrent",
		Category: NodeCategoryTorrent,
		Io: map[string]*NodeIO{
			"url": {
				Id:        "url",
				ValueType: ValueTypeString,
				Label:     "Url",
				Required:  true,
				Type:      IOTypeInput,
			},
			"directory": {
				Id:        "directory",
				ValueType: ValueTypeDirectory,
				Label:     "Directory",
				Required:  true,
				Value:     "./",
				Type:      IOTypeInput,
			},
			"service": {
				Id:        "service",
				ValueType: ValueTypeTorrentService,
				Required:  true,
				Type:      IOTypeConnectedOnly,
				Label:     "Torrent Service",
			},
		},
		Status: StatusPending,
		Run:    AddTorrentNodeFunc,
	}
}

func AddTorrentNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	log.Info().Any("nodeIo", node.Io["service"].Value).Send()
	service, ok := utils.FromAny[services.ITorrentService](node.Io["service"].Value).Value()

	if !ok || service == nil {
		return NewInvalidNodeIOError(&node, "service")
	}

	dir, ok := utils.FromAny[string](node.Io["directory"].Value).Value()
	if !ok || dir == "" {
		return NewInvalidNodeIOError(&node, "directory")
	}

	link, ok := utils.FromAny[string](node.Io["url"].Value).Value()
	if !ok || link == "" {
		return NewInvalidNodeIOError(&node, "link")
	}

	_, err := service.AddTorrent(ctx, dir, link)
	if err != nil {
		return err
	}

	return nil
}
