package models

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func NewRefreshRequest(tok string) *RefreshRequest {
	return &RefreshRequest{RefreshToken: tok}
}
