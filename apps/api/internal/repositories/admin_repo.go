package repositories

import (
	//"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"gorm.io/gorm"
)

type AdminRepo struct {
	gormdb *gorm.DB
}

func NewAdminRepo(dbclient *client.PostgresClient) *AdminRepo {
	return &AdminRepo{
		gormdb: dbclient.GetClient(),
	}
}

func (a *AdminRepo) DoesAdminExist() (bool, error) {
	var count int64
	if err := a.gormdb.Model(models.Admin{}).Count(&count).Error; err != nil {
		return false, utils.NewDatabaseError("failed to check admin existence", err)
	}
	return count > 0, nil
}

func (a *AdminRepo) CreateAdmin(admin *models.Admin) error {
	err := a.gormdb.FirstOrCreate(admin, &models.Admin{ID: 1}).Error
	if err != nil {
		return utils.NewDatabaseError("failed to create admin", err)
	}
	return nil
}
