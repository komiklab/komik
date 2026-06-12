package auth

import (
	"errors"
	"net/http"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"github.com/komiklab/komik/internal/utils"
)

type AuthService struct {
	authrepo *repositories.AdminRepo
}

func NewAuthService(dbcon *client.PostgresClient) *AuthService {
	authrepo := repositories.NewAdminRepo(dbcon)
	return &AuthService{
		authrepo: authrepo,
	}
}

func (a *AuthService) DoesAdminExist() (bool, error) {
	exist, err := a.authrepo.DoesAdminExist()
	if err != nil {
		if errors.Is(err, utils.NewDatabaseError) {
			komikErr := err.(*utils.KomikError)
			komikErr = komikErr.WithStatusCode(http.StatusInternalServerError)
			return false, komikErr
		}
		return false, err
	}
	return exist, nil
}

func (a *AuthService) CreateAdmin(admin *models.Admin) error {
	var err error
	admin.Password, err = utils.HashPassword(admin.Password)
	if err != nil {
		return utils.NewInternalServerError("failed to hash password")
	}
	return a.authrepo.CreateAdmin(admin)
}
