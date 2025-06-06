package resumeService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
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

	redisRepo, err := s.resumeRepository.NewCacheClient()
	if err != nil {
		return resume.ResumeDetailDTO{}, err
	}

	redisKey := fmt.Sprintf("cv-%s", resumeID)

	cachedResume, err := redisRepo.Get(ctx, redisKey)
	if err == nil && cachedResume != "" {
		var cachedResumeData entity.ResumeDetail
		err = json.Unmarshal([]byte(cachedResume), &cachedResumeData)
		if err != nil {
			return resume.ResumeDetailDTO{}, err
		}

		if cachedResumeData.UserID != userID {
			return resume.ResumeDetailDTO{}, resume.ErrIncorrectObjectID
		}

		return s.formattedResumeDetail(cachedResumeData), nil
	}

	resumeData, err := mongoRepo.Resume.GetByIDAndUserID(ctx, resumeID, userID)
	if err != nil {
		return resume.ResumeDetailDTO{}, err
	}

	return s.formattedResumeDetail(resumeData), nil
}

func (s resumeService) DeleteResumeByID(ctx context.Context, resumeID string) error {
	sqlRepo, err := s.resumeRepository.NewSqlClient(false)
	if err != nil {
		return err
	}

	err = sqlRepo.Resume.DeleteById(ctx, resumeID)
	if err != nil {
		return err
	}

	return nil
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

		updatedResume, err := json.Marshal(s.formattedResumeDetail(resumeData))
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

	updatedResume, err := json.Marshal(s.formattedResumeDetail(resumeData))
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
			AdviceMessage:  strings.Split(result.AdviceMessage, "\n"),
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
		AdviceMessage:  strings.Split(result.AdviceMessage, "\n"),
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
		AdviceMessage:  strings.Split(result.AdviceMessage, "\n"),
	}, nil
}

func (s resumeService) JobVacancyFromResume(ctx context.Context, resumeID string, userID string) ([]resume.JobVacancyRespone, error) {
	redisRepo, err := s.resumeRepository.NewCacheClient()
	if err != nil {
		return nil, err
	}

	mongoRepo, err := s.resumeRepository.NewMongoClient(ctx, false)
	if err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("cv-%s", resumeID)

	cachedResume, err := redisRepo.Get(ctx, redisKey)
	if err == nil && cachedResume != "" {
		var cachedResumeData resume.ResumeDetailDTO
		err = json.Unmarshal([]byte(cachedResume), &cachedResumeData)
		if err != nil {
			return nil, err
		}

		if cachedResumeData.UserID != userID {
			return nil, resume.ErrIncorrectObjectID
		}

		experienceAndSkills, err := s.google.Gemini.GenerateExperienceAndSkillsParagrafFromJSON(ctx, []byte(cachedResume))
		if err != nil {
			return nil, err
		}

		result, err := s.machineLearning.FindJobsRelated(ctx, experienceAndSkills)
		if err != nil {
			return nil, err
		}

		res := make([]resume.JobVacancyRespone, len(result))
		for i, j := range result {
			res[i] = s.formattedJobVacancy(&j)
		}

		return res, nil
	}

	resumeData, err := mongoRepo.Resume.GetByIDAndUserID(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	jsonResumeData, err := json.Marshal(s.formattedResumeDetail(resumeData))
	if err != nil {
		return nil, err
	}

	experienceAndSkills, err := s.google.Gemini.GenerateExperienceAndSkillsParagrafFromJSON(ctx, jsonResumeData)
	if err != nil {
		return nil, err
	}

	result, err := s.machineLearning.FindJobsRelated(ctx, experienceAndSkills)
	if err != nil {
		return nil, err
	}

	res := make([]resume.JobVacancyRespone, len(result))
	for i, j := range result {
		res[i] = s.formattedJobVacancy(&j)
	}

	return res, nil
}

func (s resumeService) JobVacancyFromPDF(ctx context.Context, pdf *multipart.FileHeader, userID string) ([]resume.JobVacancyRespone, error) {
	pdfByte, err := pdf.Open()
	if err != nil {
		return nil, err
	}

	defer func(pdfByte multipart.File) {
		err := pdfByte.Close()
		if err != nil {
		}
	}(pdfByte)

	fileBytes, err := io.ReadAll(pdfByte)
	if err != nil {
		return nil, err
	}

	experienceAndSkills, err := s.google.Gemini.GenerateExperienceAndSkillsParagrafFromPDF(ctx, fileBytes)
	if err != nil {
		return nil, err
	}

	if experienceAndSkills == "not-resume" {
		return nil, resume.ErrInvalidResumeFile
	}

	result, err := s.machineLearning.FindJobsRelated(ctx, experienceAndSkills)
	if err != nil {
		return nil, err
	}

	res := make([]resume.JobVacancyRespone, len(result))
	for i, j := range result {
		res[i] = s.formattedJobVacancy(&j)
	}

	return res, nil
}

func (s resumeService) DownloadResumePDF(ctx context.Context, resumeID string, userID string) ([]byte, error) {
	redisRepo, err := s.resumeRepository.NewCacheClient()
	if err != nil {
		return nil, err
	}

	mongoRepo, err := s.resumeRepository.NewMongoClient(ctx, false)
	if err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("cv-%s", resumeID)
	cachedResume, err := redisRepo.Get(ctx, redisKey)
	if err == nil && cachedResume != "" {
		var cachedResumeData resume.ResumeDetailDTO
		err = json.Unmarshal([]byte(cachedResume), &cachedResumeData)
		if err != nil {
			return nil, err
		}

		if cachedResumeData.UserID != userID {
			return nil, resume.ErrIncorrectObjectID
		}

		tmpl, err := template.ParseFiles("resume_template.gohtml")
		if err != nil {
			return nil, err
		}

		var htmlBuffer bytes.Buffer
		err = tmpl.Execute(&htmlBuffer, cachedResumeData)
		if err != nil {
			return nil, err
		}

		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			return nil, err
		}

		pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlBuffer.String())))

		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

		err = pdfg.Create()
		if err != nil {
			return nil, err
		}

		return pdfg.Bytes(), nil
	}

	resumeData, err := mongoRepo.Resume.GetByIDAndUserID(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles("resume_template.gohtml")
	if err != nil {
		return nil, err
	}

	// Create a buffer to hold the rendered HTML
	var htmlBuffer bytes.Buffer
	formatted := s.formattedResumeDetail(resumeData)
	err = tmpl.Execute(&htmlBuffer, formatted)
	if err != nil {
		return nil, err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlBuffer.String())))
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
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
