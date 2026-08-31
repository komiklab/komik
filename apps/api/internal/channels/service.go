package channels

import (
	"context"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

type ChannelService struct {
	channelRepo *repositories.ChannelRepo
}

func NewChannelService(dbcon *client.PostgresClient) *ChannelService {
	return &ChannelService{
		channelRepo: repositories.NewChannelRepo(dbcon),
	}
}

func (s *ChannelService) List(ctx context.Context) (*models.Channellist, error) {
	channels, err := s.channelRepo.FetchChannels()
	if err != nil {
		return nil, err
	}
	return &models.Channellist{Channels: channels}, nil
}
