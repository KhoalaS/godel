package pipeline

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hekmon/transmissionrpc/v3"
)

type TransmissionTorrentService struct {
	serverUrl string
	client    *transmissionrpc.Client
}

func NewTransmissionTorrentService(config TransmissionServerConfig) (*TransmissionTorrentService, error) {
	serverUrl, err := url.Parse(fmt.Sprintf("http://%s:%s@%s", config.Username, config.Password, config.ServerUrl))

	if err != nil {
		return nil, err
	}

	client, err := transmissionrpc.New(serverUrl, nil)
	if err != nil {
		return nil, err
	}

	return &TransmissionTorrentService{
		serverUrl: config.ServerUrl,
		client:    client,
	}, nil
}

func (t *TransmissionTorrentService) AddTorrent(ctx context.Context, downloadDirectory string, magnetLink string) (string, error) {

	torrent, err := t.client.TorrentAdd(ctx, transmissionrpc.TorrentAddPayload{
		DownloadDir: &downloadDirectory,
		Filename:    &magnetLink,
	})

	if err != nil {
		return "", err
	}

	id := strconv.Itoa(int(*torrent.ID))

	return id, nil
}

func (t *TransmissionTorrentService) PauseTorrent(ctx context.Context, id string) error {
	return nil
}
func (t *TransmissionTorrentService) RemoveTorrent(ctx context.Context, id string) error {
	return nil
}

type TransmissionServerConfig struct {
	Username  string
	Password  string
	ServerUrl string
}
