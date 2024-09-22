package authService

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	authRepository "github.com/bccfilkom/career-path-service/internal/api/authentication/repository"
	"golang.org/x/net/context"
)

type authService struct {
	authRepository authRepository.Repository
}

type AuthService interface {
	RegisterUser(ctx context.Context, req authentication.CreateUserRequest) error
}

func New(authRepository authRepository.Repository) AuthService {
	return &authService{
		authRepository: authRepository,
	}
}
