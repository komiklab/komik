package channels

import (
	"context"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
)

func GetChannelList(ctx context.Context,dbcon *client.PostgresClient) (*models.Channellist, error) {
	s := NewChannelService(dbcon)
	return s.List(ctx)
}