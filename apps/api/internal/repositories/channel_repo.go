package repositories

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"gorm.io/gorm"
)

type ChannelRepo struct {
	gormdb *gorm.DB
}

func (repo *ChannelRepo) GetChannelByID(channelId string) (*models.Channel, error) {
	var channel models.Channel
	err := repo.gormdb.First(&channel, "id = ?", channelId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFoundError("channel not found", err)
		}
		return nil, utils.NewDatabaseError("failed to get channel by ID", err)
	}
	return &channel, nil
}

func NewChannelRepo(dbclient *client.PostgresClient) *ChannelRepo {
	return &ChannelRepo{
		gormdb: dbclient.GetClient(),
	}
}

func (repo *ChannelRepo) FetchChannels() ([]models.Channel, error) {
	var channels []models.Channel
	err := repo.gormdb.Find(&channels).Error
	if err != nil {
		return nil, utils.NewDatabaseError("failed to fetch channels", err)
	}
	return channels, nil
}

func (repo *ChannelRepo) CreateChannel(channel *models.Channel) error {
	err := repo.gormdb.Create(channel).Error
	if err != nil {
		return utils.NewDatabaseError("failed to create channel", err)
	}
	return nil
}
