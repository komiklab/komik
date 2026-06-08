package auth

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"github.com/komiklab/komik/internal/utils"
)

type AuthService struct {
	authrepo *repositories.AdminRepo
}

func NewAuthService(dbcon client.Client) *AuthService {
	authrepo := repositories.NewAdminRepo(dbcon)
	return &AuthService{
		authrepo: authrepo,
	}
}

func (a *AuthService) DoesAdminExist() (bool, error) {
	return a.authrepo.DoesAdminExist()
}

func (a *AuthService) CreateAdmin(admin *models.Admin) error {
	var err error
	admin.Password, err = utils.HashPassword(admin.Password)
	if err != nil {
		return err
	}
	return a.authrepo.CreateAdmin(admin)
}
