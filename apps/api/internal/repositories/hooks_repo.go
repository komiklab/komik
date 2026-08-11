package repositories

import (
	//"errors"

	//"github.com/google/uuid"
	"errors"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HooksRepo struct {
	gormdb *gorm.DB
}

func NewHooksRepo(dbclient *client.PostgresClient) *HooksRepo {
	return &HooksRepo{
		gormdb: dbclient.GetClient(),
	}
}

func (h *HooksRepo) CreateHook(hook *models.Hooks) error {
	// create hook if only name is unique
	result := h.gormdb.Clauses(clause.OnConflict{DoNothing: true}).Create(hook)
	if result.RowsAffected == 0 {
		return utils.NewConflictError("hook already exist", errors.New("hook already exist"))
	}
	return nil
}

// func (h *HooksRepo) DeleteHook(id uuid.UUID) error {
// 	err := h.gormdb.Where("id = ?", id).Delete(&models.Hooks{}).Error
// 	if err != nil {
// 		return utils.NewDatabaseError("failed to delete hook", err)
// 	}
// 	return nil
// }

func (h *HooksRepo) FetchHook(name string) (*models.Hooks, error) {
	var hook models.Hooks
	err := h.gormdb.Where("name = ?", name).First(&hook).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAuthenticationError("hook not found", err)
		}
		return nil, utils.NewDatabaseError("failed to fetch hook", err)
	}
	return &hook, nil
}

func (h *HooksRepo) FetchHooks() ([]models.Hooks, error) {
	var hooks []models.Hooks
	err := h.gormdb.Find(&hooks).Error
	if err != nil {
		return nil, utils.NewDatabaseError("failed to fetch hooks", err)
	}
	return hooks, nil
}
