package resumeRepository

import (
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func (r *resumeRepository) CreateResume(ctx context.Context, resume entity.ResumeDetail) (*mongo.InsertOneResult, error) {
	return r.db.Collection("resume").InsertOne(ctx, resume)
}

func (r *resumeRepository) GetByIDAndUserID(ctx context.Context, ID string, userID string) (entity.ResumeDetail, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return entity.ResumeDetail{}, resume.ErrIncorrectObjectID
	}

	var result entity.ResumeDetail
	err = r.db.Collection("resume").FindOne(ctx, bson.M{"_id": objectID, "userID": userID}).Decode(&result)
	if err != nil {
		return entity.ResumeDetail{}, err
	}

	return result, nil
}
