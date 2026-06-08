package repositories

import (
	//"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"gorm.io/gorm"
)

type AdminRepo struct {
	gormdb *gorm.DB
}

func NewAdminRepo(dbclient client.Client) *AdminRepo {
	return &AdminRepo{
		gormdb: dbclient.GetClient().(*gorm.DB),
	}
}

func (a *AdminRepo) DoesAdminExist() (bool, error) {
	var count int64
	if err := a.gormdb.Model(models.Admin{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a *AdminRepo) CreateAdmin(admin *models.Admin) error {
	return a.gormdb.FirstOrCreate(admin, &models.Admin{ID: 1}).Error
}
