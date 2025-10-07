package pipeline

import (
	"context"
	"fmt"
	"os"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/rs/zerolog/log"
)

func CreateTransmissionNode() Node {
	return Node{
		Type:     "transmission-service",
		Name:     "Transmission Service",
		Category: NodeCategoryTorrent,
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

	ok, serverVersion, serverMinimumVersion, err := transmissionTorrentService.client.RPCVersion(ctx)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
	}
	log.Info().Msg(fmt.Sprintf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)\n",
		serverVersion, transmissionrpc.RPCVersion))

	node.Io["service"].Value = transmissionTorrentService
	return nil
}
