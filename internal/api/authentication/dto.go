package authentication

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Paassword string `json:"password" validate:"required,min=8,max=32"`
	Name      string `json:"name" validate:"required,min=3,max=255"`
}
