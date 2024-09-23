package authentication

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Paassword string `json:"password" validate:"required,min=8,max=32"`
	Name      string `json:"name" validate:"required,min=3,max=255"`
}

type SigninUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SinginUserResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresInSeconds int64  `json:"expiresInSeconds"`
	SessionID        string `json:"sessionID"`
}

type RefresTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken      string `json:"accessToken"`
	ExpiresInSeconds int64  `json:"expiresInSeconds"`
}
