package dto

type DeveloperDTO struct {
	UserID    uint   `json:"user_id"`
	Handle    string `json:"handle"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"email"`
	UserRole  string `gorm:"string"`
	IsBanned  bool   `gorm:"is_banned"`
}
