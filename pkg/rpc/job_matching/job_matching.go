package jobMatching

import (
	"context"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type machineLearning struct {
	client CareerPathClient
}

type MachineLearning interface {
	FindJobsRelated(ctx context.Context, experienceAndSkills string) ([]Job, error)
	ResumeScoring(ctx context.Context, resumeString string) (CompareCVResponse, error)
}

func NewRpcClient(log *logrus.Logger) MachineLearning {
	serverAddr := env.GetString("MACHINE_LEARNING_ENDPOINT", "")

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("grpc.NewClient job_matching err: %v", err)
	}

	client := NewCareerPathClient(conn)

	return &machineLearning{
		client: client,
	}
}

func (m machineLearning) FindJobsRelated(ctx context.Context, experienceAndSkills string) ([]Job, error) {
	req := FindJobsRequest{
		ExperienceAndSkills: experienceAndSkills,
	}

	jobs, err := m.client.FindJobs(ctx, &req)
	if err != nil {
		return nil, err
	}

	res := make([]Job, len(jobs.Jobs))
	for i, j := range jobs.Jobs {
		res[i] = Job{
			Company:     j.Company,
			Link:        j.Link,
			Title:       j.Title,
			Score:       j.Score,
			Description: j.Description,
		}
	}

	return res, nil
}

func (m machineLearning) ResumeScoring(ctx context.Context, resumeString string) (CompareCVResponse, error) {
	req := CompareCVRequest{
		CvJson: resumeString,
	}

	resultScore, err := m.client.CompareCV(ctx, &req)
	if err != nil {
		return CompareCVResponse{}, err
	}

	return CompareCVResponse{
		FinalScore:     resultScore.FinalScore,
		AdviceMessage:  resultScore.AdviceMessage,
		OverallMessage: resultScore.OverallMessage,
	}, nil
}
