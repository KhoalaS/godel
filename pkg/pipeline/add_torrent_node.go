package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/services"
	"github.com/pkg/errors"
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
	service, ok := node.Io["service"].Value.(services.ITorrentService)

	if !ok || service == nil {
		return errors.New("no torrent service provided")
	}

	dir, ok := node.Io["directory"].Value.(string)
	if !ok || dir == "" {
		return errors.New("no directory provided for torrent")
	}

	link, ok := node.Io["url"].Value.(string)
	if !ok || link == "" {
		return errors.New("no link provided for torrent")
	}

	_, err := service.AddTorrent(ctx, dir, link)
	if err != nil {
		return err
	}

	return nil
}
