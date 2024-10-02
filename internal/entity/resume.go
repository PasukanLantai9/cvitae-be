package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Resume struct {
	ID        string
	Name      string
	UserID    string
	CreatedAt time.Time
}

type ResumeDetail struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	UserID                 string             `bson:"userID,omitempty"`
	PersonalDetails        PersonalDetails    `bson:"personalDetails,omitempty"`
	ProfessionalExperience []Experience       `bson:"professionalExperience,omitempty"`
	Education              []Education        `bson:"education,omitempty"`
	LeadershipExperience   []Leadership       `bson:"leadershipExperience,omitempty"`
	Others                 []Achievement      `bson:"others,omitempty"`
}

type PersonalDetails struct {
	FullName      string `bson:"fullName,omitempty"`
	PhoneNumber   string `bson:"phoneNumber,omitempty"`
	Email         string `bson:"email,omitempty"`
	Linkedin      string `bson:"linkedin,omitempty"`
	PortfolioURL  string `bson:"portfolioUrl,omitempty"`
	Description   string `bson:"description,omitempty"`
	AddressString string `bson:"addressString,omitempty"`
}

type Experience struct {
	StartDate   Date          `bson:"startDate,omitempty"`
	EndDate     Date          `bson:"endDate,omitempty"`
	RoleTitle   string        `bson:"roleTitle,omitempty"`
	CompanyName string        `bson:"companyName,omitempty"`
	Location    string        `bson:"location,omitempty"`
	Current     bool          `bson:"current,omitempty"`
	Elaboration []Elaboration `bson:"elaboration,omitempty"`
}

type Education struct {
	StartDate   Date          `bson:"startDate,omitempty"`
	EndDate     Date          `bson:"endDate,omitempty"`
	School      string        `bson:"school,omitempty"`
	Location    string        `bson:"location,omitempty"`
	DegreeLevel string        `bson:"degreeLevel,omitempty"`
	Major       string        `bson:"major,omitempty"`
	GPA         float64       `bson:"gpa,omitempty"`
	MaxGPA      float64       `bson:"maxGpa,omitempty"`
	Elaboration []Elaboration `bson:"elaboration,omitempty"`
}

type Leadership struct {
	StartDate        Date          `bson:"startDate,omitempty"`
	EndDate          Date          `bson:"endDate,omitempty"`
	RoleTitle        string        `bson:"roleTitle,omitempty"`
	OrganisationName string        `bson:"organisationName,omitempty"`
	Location         string        `bson:"location,omitempty"`
	Current          bool          `bson:"current,omitempty"`
	Elaboration      []Elaboration `bson:"elaboration,omitempty"`
}

type Achievement struct {
	Elaboration Elaboration `bson:"elaboration,omitempty"`
	Name        string      `bson:"name,omitempty"`
	Date        Date        `bson:"date,omitempty"`
	Category    string      `bson:"category,omitempty"`
}

type Date struct {
	Month string `bson:"month,omitempty"`
	Year  int    `bson:"year,omitempty"`
}

type Elaboration struct {
	Text string `bson:"text,omitempty"`
}
