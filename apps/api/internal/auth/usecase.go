package auth


type AuthService struct {
	
}

func NewAuthService() *AuthService {
	return &AuthService{}
}


func (a *AuthService) DoesAdminExist() (bool, error) {
	return false, nil
}

func (a *AuthService) CreateAdmin() error {
	return nil
}