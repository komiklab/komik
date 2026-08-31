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

func CreateChannel(ctx context.Context, dbcon *client.PostgresClient, channel *models.Channel) error {
	s := NewChannelService(dbcon)
	return s.Create(ctx, channel)
}

func GetChannelByID(ctx context.Context, dbcon *client.PostgresClient, channelId string) (*models.Channel, error) {
	s := NewChannelService(dbcon)
	return s.GetByID(ctx, channelId)
}