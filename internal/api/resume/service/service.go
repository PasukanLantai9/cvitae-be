package resumeService

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	resumeRepository "github.com/bccfilkom/career-path-service/internal/api/resume/repository"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type resumeService struct {
	resumeRepository resumeRepository.Repository
}

type ResumeService interface {
	CreateResume(context.Context, resume.ResumeRequest, authentication.UserClaims) (string, error)
	GetUserResume(context.Context, string) ([]resume.GetResumeResponse, error)
	GetResumeByID(context.Context, string, string) (resume.ResumeDetailDTO, error)
	UpdateResumeByID(context.Context, entity.ResumeDetail) error

	SyncResumesFromRedisToMongo(redisClient *redis.Client, cvCollection *mongo.Collection) error
}

func New(repo resumeRepository.Repository) ResumeService {
	return &resumeService{
		resumeRepository: repo,
	}
}
