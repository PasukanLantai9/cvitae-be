package resumeService

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"time"
)

func (s resumeService) CreateResume(ctx context.Context, req resume.ResumeRequest, user authentication.UserClaims) (string, error) {
	sqlRepo, err := s.resumeRepository.NewSqlClient(true)
	if err != nil {
		return "", err
	}

	mongoRepo, err := s.resumeRepository.NewMongoClient(ctx, true)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			if errTx := sqlRepo.Rollback(); errTx != nil {
				err = errTx
				return
			}

			if errTx := mongoRepo.Rollback(); errTx != nil {
				err = errTx
				return
			}
		}
	}()

	resumeReq := entity.Resume{
		UserID:    user.ID,
		CreatedAt: time.Now(),
		Name:      req.Name,
	}

	resumeDetail := entity.ResumeDetail{
		UserID:                 user.ID,
		Education:              []entity.Education{{}},
		LeadershipExperience:   []entity.Leadership{{}},
		ProfessionalExperience: []entity.Experience{{}},
		Others:                 []entity.Achievement{{}},
	}

	result, err := mongoRepo.Resume.CreateResume(ctx, resumeDetail)
	if err != nil {
		return "", err
	}

	resumeReq.ID = result.InsertedID.(primitive.ObjectID).Hex()
	if err := sqlRepo.Resume.CreateResume(ctx, resumeReq); err != nil {
		return "", err
	}

	if err := mongoRepo.Commit(); err != nil {
		return "", err
	}

	if err := sqlRepo.Commit(); err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s resumeService) GetUserResume(ctx context.Context, userID string) ([]resume.GetResumeResponse, error) {
	sqlRepo, err := s.resumeRepository.NewSqlClient(false)
	if err != nil {
		return nil, err
	}

	resumes, err := sqlRepo.Resume.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resumesFormatted := make([]resume.GetResumeResponse, len(resumes))
	for i, resumeData := range resumes {
		resumesFormatted[i] = s.formattedResume(resumeData)
	}

	return resumesFormatted, nil
}

func (s resumeService) GetResumeByID(ctx context.Context, resumeID string, userID string) (resume.ResumeDetailDTO, error) {
	mongoRepo, err := s.resumeRepository.NewMongoClient(ctx, false)
	if err != nil {
		return resume.ResumeDetailDTO{}, err
	}

	resumeData, err := mongoRepo.Resume.GetByIDAndUserID(ctx, resumeID, userID)
	if err != nil {
		return resume.ResumeDetailDTO{}, err
	}

	return s.formattedResumeDetail(resumeData), nil
}
