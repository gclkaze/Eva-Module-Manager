package services

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (inst *AuthService) Register(email string, pwd string) (string, error) {
	return "", nil
}

func (inst *AuthService) Login(email string, pwd string) (string, error) {
	return "", nil
}

func (inst *AuthService) Logout(email string) (string, error) {
	return "", nil
}

/*func (inst *AuthService) SetActiveUser(email string) (string, error) {
	return "", nil
}*/

func (inst *AuthService) SwitchActiveUser(email string) (string, error) {
	return "", nil
}

func (inst *AuthService) GetStatus() error {
	return nil
}
