package resumeService

import (
	"encoding/json"
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"io"
	"mime/multipart"
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
		PersonalDetails:        entity.PersonalDetails{},
		Education:              make([]entity.Education, 0),
		LeadershipExperience:   make([]entity.Leadership, 0),
		ProfessionalExperience: make([]entity.Experience, 0),
		Others:                 make([]entity.Achievement, 0),
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

func (s resumeService) UpdateResumeByID(ctx context.Context, resumeData entity.ResumeDetail) error {
	redisRepo, err := s.resumeRepository.NewCacheClient()
	if err != nil {
		return err
	}

	mongoRepo, err := s.resumeRepository.NewMongoClient(ctx, false)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("cv-%s", resumeData.ID.Hex())

	cachedResume, err := redisRepo.Get(ctx, redisKey)
	if err == nil && cachedResume != "" {
		var cachedResumeData resume.ResumeDetailDTO
		err = json.Unmarshal([]byte(cachedResume), &cachedResumeData)
		if err != nil {
			return err
		}

		if cachedResumeData.UserID != resumeData.UserID {
			return resume.ErrIncorrectObjectID
		}

		updatedResume, err := json.Marshal(resumeData)
		if err != nil {
			return err
		}

		err = redisRepo.Set(ctx, redisKey, updatedResume, 0)
		if err != nil {
			return err
		}

		return nil
	}

	_, err = mongoRepo.Resume.GetByIDAndUserID(ctx, resumeData.ID.Hex(), resumeData.UserID)
	if err != nil {
		return err
	}

	err = mongoRepo.Resume.Update(ctx, resumeData)
	if err != nil {
		return err
	}

	updatedResume, err := json.Marshal(resumeData)
	if err != nil {
		return err
	}

	err = redisRepo.Set(ctx, redisKey, updatedResume, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s resumeService) ScoringResume(ctx context.Context, resumeID string, userID string) (resume.ScoringResumeResponse, error) {
	redisRepo, err := s.resumeRepository.NewCacheClient()
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	mongoRepo, err := s.resumeRepository.NewMongoClient(ctx, false)
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	redisKey := fmt.Sprintf("cv-%s", resumeID)

	cachedResume, err := redisRepo.Get(ctx, redisKey)
	if err == nil && cachedResume != "" {
		var cachedResumeData resume.ResumeDetailDTO
		err = json.Unmarshal([]byte(cachedResume), &cachedResumeData)
		if err != nil {
			return resume.ScoringResumeResponse{}, err
		}

		if cachedResumeData.UserID != userID {
			return resume.ScoringResumeResponse{}, resume.ErrIncorrectObjectID
		}

		result, err := s.machineLearning.ResumeScoring(ctx, cachedResume)
		if err != nil {
			return resume.ScoringResumeResponse{}, err
		}

		return resume.ScoringResumeResponse{
			OverallMessage: result.OverallMessage,
			FinalScore:     float64(result.FinalScore),
			AdviceMessage:  result.AdviceMessage,
		}, nil
	}

	resumeData, err := mongoRepo.Resume.GetByIDAndUserID(ctx, resumeID, userID)
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	jsonResumeData, err := json.Marshal(s.formattedResumeDetail(resumeData))
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	jsonString := string(jsonResumeData)

	result, err := s.machineLearning.ResumeScoring(ctx, jsonString)
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	return resume.ScoringResumeResponse{
		OverallMessage: result.OverallMessage,
		FinalScore:     float64(result.FinalScore),
		AdviceMessage:  result.AdviceMessage,
	}, nil
}

func (s resumeService) ScoringResumePDF(ctx context.Context, pdfFile *multipart.FileHeader, userID string) (resume.ScoringResumeResponse, error) {
	pdfByte, err := pdfFile.Open()
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}
	defer func(pdfByte multipart.File) {
		err := pdfByte.Close()
		if err != nil {

		}
	}(pdfByte)

	fileBytes, err := io.ReadAll(pdfByte)
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	jsonString, err := s.google.Gemini.GenerateResumeJsonFromPDF(ctx, fileBytes)
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	if jsonString == "not-resume" {
		return resume.ScoringResumeResponse{}, resume.ErrInvalidResumeFile
	}

	result, err := s.machineLearning.ResumeScoring(ctx, jsonString)
	if err != nil {
		return resume.ScoringResumeResponse{}, err
	}

	return resume.ScoringResumeResponse{
		OverallMessage: result.OverallMessage,
		FinalScore:     float64(result.FinalScore),
		AdviceMessage:  result.AdviceMessage,
	}, nil
}

func (s resumeService) SyncResumesFromRedisToMongo(redisClient *redis.Client, cvCollection *mongo.Collection) error {
	ctx := context.Background()

	keys, err := redisClient.Keys(ctx, "cv-*").Result()
	if err != nil {
		fmt.Println("Error getting keys from Redis:", err)
		return err
	}

	for _, key := range keys {
		cvData, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			fmt.Println("Error getting CV from Redis for key", key, ":", err)
			continue
		}

		var cv entity.ResumeDetail
		err = json.Unmarshal([]byte(cvData), &cv)
		if err != nil {
			fmt.Println("Error unmarshalling CV data for key", key, ":", err)
			continue
		}

		objectID, _ := primitive.ObjectIDFromHex(cv.ID.Hex())

		_, err = cvCollection.UpdateOne(ctx, bson.M{"userID": cv.UserID, "_id": objectID}, bson.M{
			"$set": bson.M{
				"personalDetails":        cv.PersonalDetails,
				"professionalExperience": cv.ProfessionalExperience,
				"education":              cv.Education,
				"leadershipExperience":   cv.LeadershipExperience,
				"others":                 cv.Others,
			},
		}, options.Update().SetUpsert(true))
		if err != nil {
			fmt.Println("Error saving CV to MongoDB for userID", cv.UserID, ":", err)
			continue
		}

		err = redisClient.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("Error deleting CV from Redis for key", key, ":", err)
		}
	}

	return nil
}
