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

func NewChannelRepo(dbclient *client.PostgresClient) *ChannelRepo {
	return &ChannelRepo{
		gormdb: dbclient.GetClient(),
	}
}

func (repo *ChannelRepo) FetchChannels()([]models.Channel,error){
	var channels []models.Channel
	err := repo.gormdb.Find(&channels).Error
	if err != nil {
		return nil, utils.NewDatabaseError("failed to fetch channels", err)
	}
	return channels, nil
}

