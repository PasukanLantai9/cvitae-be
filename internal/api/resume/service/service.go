package resumeService

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	resumeRepository "github.com/bccfilkom/career-path-service/internal/api/resume/repository"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/pkg/rpc/job_matching"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type resumeService struct {
	resumeRepository resumeRepository.Repository
	machineLearning  jobMatching.MachineLearning
}

type ResumeService interface {
	CreateResume(context.Context, resume.ResumeRequest, authentication.UserClaims) (string, error)
	GetUserResume(context.Context, string) ([]resume.GetResumeResponse, error)
	GetResumeByID(context.Context, string, string) (resume.ResumeDetailDTO, error)
	UpdateResumeByID(context.Context, entity.ResumeDetail) error

	ScoringResume(context.Context, string, string) (resume.ScoringResumeResponse, error)

	SyncResumesFromRedisToMongo(redisClient *redis.Client, cvCollection *mongo.Collection) error
}

func New(repo resumeRepository.Repository, ml jobMatching.MachineLearning) ResumeService {
	return &resumeService{
		resumeRepository: repo,
		machineLearning:  ml,
	}
}
