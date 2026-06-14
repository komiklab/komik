package auth

import (
	"errors"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"github.com/komiklab/komik/internal/utils"
)

type AuthService struct {
	authrepo *repositories.AdminRepo
}

func NewAuthService(dbcon *client.PostgresClient, cache *client.RedisClient) *AuthService {
	authrepo := repositories.NewAdminRepo(dbcon, cache)
	return &AuthService{
		authrepo: authrepo,
	}
}

func (a *AuthService) DoesAdminExist() (bool, error) {
	exist, err := a.authrepo.DoesAdminExist()
	if err != nil {
		return false, utils.NewGeneralError(err)
	}
	return exist, nil
}

func (a *AuthService) CreateAdmin(admin *models.Admin) error {
	var err error
	admin.Password, err = utils.HashPassword(admin.Password)
	if err != nil {
		return utils.NewGeneralError(err)
	}
	return a.authrepo.CreateAdmin(admin)
}

func (a *AuthService) VerifyPassword(admin *models.Admin) error {
	passwordReceived := admin.Password
	hashedPassword, err := a.authrepo.FetchPassword(admin)
	if utils.IsErrNotNil(err) {
		return err
	}
	verified, err := utils.VerifyPassword(passwordReceived,hashedPassword)
	if utils.IsErrNotNil(err) {
		return utils.NewInternalServerError("could not verify password because of internal issue", err)
	}
	if !verified {
		return utils.NewAuthenticationError("invalid password", errors.New("wrong password"))
	}
	return nil
}

func (a *AuthService) CreateSession(user models.UserRepresentation) (string, error){
	// we will first store the user if its new in database
	err := a.authrepo.SaveUserIfNotExist(user)
	if utils.IsErrNotNil(err) {
		return "", err
	}
	// then we will create session for the user
	sessionID, err := a.authrepo.CreateSession(&user)
	if utils.IsErrNotNil(err) {
		return "", err
	}	
	return sessionID, nil
}
	

