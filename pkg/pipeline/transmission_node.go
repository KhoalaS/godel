package pipeline

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
)

func CreateTransmissionNode() Node {
	return Node{
		Type:     "transmission-service",
		Name:     "Transmission Service",
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"service": {
				Id:        "service",
				ValueType: ValueTypeTorrentService,
				Label:     "Service",
				Type:      IOTypeGenerated,
			},
		},
		Status: StatusPending,
		Run:    TransmissionNodeFunc,
	}
}

func TransmissionNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	transmissionUser := os.Getenv("TR_USERNAME")
	transmissionPassword := os.Getenv("TR_PASSWORD")
	transmissionServer := os.Getenv("TR_SERVER_URL")

	transmissionTorrentService, err := NewTransmissionTorrentService(TransmissionServerConfig{
		Username:  transmissionUser,
		Password:  transmissionPassword,
		ServerUrl: transmissionServer,
	})

	if err != nil {
		log.Err(err).Str("serverUrl", transmissionServer).Msg("could not create Transmission torrent service.")
		return err
	}

	node.Io["service"].Value = transmissionTorrentService
	return nil
}
