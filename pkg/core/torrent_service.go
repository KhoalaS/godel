package core

import (
	"context"
)

type ITOrrentService interface {
	AddTorrent(ctx context.Context, downloadDirectory string, torrentLink string) (string, error)
	PauseTorrent(ctx context.Context, id string) error
	RemoveTorrent(ctx context.Context, id string) error
}
