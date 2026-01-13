package userinput

type LoginCreds struct {
	Email    string
	Password string
}

func NewLoginCreds() *LoginCreds {
	return &LoginCreds{}
}
