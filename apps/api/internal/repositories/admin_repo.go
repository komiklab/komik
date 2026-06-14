package repositories

import (
	//"github.com/komiklab/komik/apidefn"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdminRepo struct {
	gormdb *gorm.DB
	cache  *client.RedisClient
}

func (a *AdminRepo) CreateSession(representation *models.UserRepresentation) (string, error) {
	session_id := uuid.New().String()
	sessoionBytes, err := json.Marshal(representation)
	if err != nil {
		return "", utils.NewInternalServerError("failed to create session", err)
	}
	err = a.cache.Set("session:"+session_id, sessoionBytes, utils.SESSION_TTL)
	if err != nil {
		return "", utils.NewRedisError("failed to create session", err)
	}
	return session_id, nil
}

func (a *AdminRepo) SaveUserIfNotExist(user models.UserRepresentation) error {
	err := a.gormdb.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
	if err != nil {
		return utils.NewDatabaseError("failed to save user", err)
	}
	return nil
}

func NewAdminRepo(dbclient *client.PostgresClient, cache *client.RedisClient) *AdminRepo {
	return &AdminRepo{
		gormdb: dbclient.GetClient(),
		cache:  cache,
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
	result := a.gormdb.FirstOrCreate(admin, &models.Admin{ID: 1})
	err := result.Error
	if err != nil {
		return utils.NewDatabaseError("failed to create admin", err)
	}
	if result.RowsAffected == 0 {
		return utils.NewConflictError("admin already exist", errors.New("admin already exist"))
	}
	return nil
}

func (a *AdminRepo) FetchPassword(admin *models.Admin) (string, error) {
	result := a.gormdb.Select("password").Where("username = ?", admin.Username).First(&admin)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utils.NewAuthenticationError("admin not found", err)
		}
		return "", utils.NewDatabaseError("failed to fetch password", err)
	}
	return admin.Password, nil
}
