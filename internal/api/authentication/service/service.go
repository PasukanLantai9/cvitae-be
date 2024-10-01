package authService

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	authRepository "github.com/bccfilkom/career-path-service/internal/api/authentication/repository"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/pkg/google"
	"golang.org/x/net/context"
)

type authService struct {
	authRepository authRepository.Repository
	google         google.Google
}

type AuthService interface {
	RegisterUser(ctx context.Context, req authentication.CreateUserRequest) error
	SinginUser(ctx context.Context, req authentication.SigninUserRequest) (authentication.SinginUserResponse, error)
	GenerateNewAccessToken(ctx context.Context, session entity.Session) (authentication.RefreshTokenResponse, error)
	CallbackGoogleOauth(ctx context.Context, code string) (authentication.SinginUserResponse, error)
}

func New(authRepository authRepository.Repository, googles google.Google) AuthService {
	return &authService{
		authRepository: authRepository,
		google:         googles,
	}
}
