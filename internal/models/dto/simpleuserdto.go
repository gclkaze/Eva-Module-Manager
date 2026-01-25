package dto

type SimpleUserAccountDTO struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
	IsBanned bool   `json:"is_banned"`
}
