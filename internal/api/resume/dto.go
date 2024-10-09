package resume

type ResumeRequest struct {
	Name string `json:"name" validate:"required"`
}

type GetResumeResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type ResumeDetailDTO struct {
	ID                     string          `json:"id"`
	UserID                 string          `json:"userID"`
	PersonalDetails        PersonalDetails `json:"personalDetails"`
	ProfessionalExperience []Experience    `json:"professionalExperience"`
	Education              []Education     `json:"education"`
	LeadershipExperience   []Leadership    `json:"leadershipExperience"`
	Others                 []Achievement   `json:"others"`
}

type PersonalDetails struct {
	FullName      string `json:"fullName,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	Email         string `json:"email,omitempty"`
	Linkedin      string `json:"linkedin,omitempty"`
	PortfolioURL  string `json:"portfolioUrl,omitempty"`
	Description   string `json:"description,omitempty"`
	AddressString string `json:"addressString,omitempty"`
}

type Experience struct {
	StartDate   Date          `json:"startDate,omitempty"`
	EndDate     Date          `json:"endDate,omitempty"`
	RoleTitle   string        `json:"roleTitle,omitempty"`
	CompanyName string        `json:"companyName,omitempty"`
	Location    string        `json:"location,omitempty"`
	Current     bool          `json:"current,omitempty"`
	Elaboration []Elaboration `json:"elaboration,omitempty"`
}

type Education struct {
	StartDate   Date          `json:"startDate,omitempty"`
	EndDate     Date          `json:"endDate,omitempty"`
	School      string        `json:"school,omitempty"`
	Location    string        `json:"location,omitempty"`
	DegreeLevel string        `json:"degreeLevel,omitempty"`
	Major       string        `json:"major,omitempty"`
	GPA         float64       `json:"gpa,omitempty"`
	MaxGPA      float64       `json:"maxGpa,omitempty"`
	Elaboration []Elaboration `json:"elaboration,omitempty"`
}

type Leadership struct {
	StartDate        Date          `json:"startDate,omitempty"`
	EndDate          Date          `json:"endDate,omitempty"`
	RoleTitle        string        `json:"roleTitle,omitempty"`
	OrganisationName string        `json:"organisationName,omitempty"`
	Location         string        `json:"location,omitempty"`
	Current          bool          `json:"current,omitempty"`
	Elaboration      []Elaboration `json:"elaboration,omitempty"`
}

type Achievement struct {
	Elaboration Elaboration `json:"elaboration,omitempty"`
	Name        string      `json:"name,omitempty"`
	Date        Date        `json:"date,omitempty"`
	Category    string      `json:"category,omitempty"`
}

type Date struct {
	Month string `json:"month,omitempty"`
	Year  int    `json:"year,omitempty"`
}

type Elaboration struct {
	Text string `json:"text,omitempty"`
}

type ScoringResumeResponse struct {
	FinalScore     float64 `json:"finalScore"`
	AdviceMessage  string  `json:"adviceMessage"`
	OverallMessage string  `json:"overallMessage"`
}
