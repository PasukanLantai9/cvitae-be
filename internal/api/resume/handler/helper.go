package resumeHandler

import (
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

func (h *ResumeHandler) convertDTOToEntity(req resume.ResumeDetailDTO) (entity.ResumeDetail, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return entity.ResumeDetail{}, err
	}

	result := entity.ResumeDetail{
		ID:     id,
		UserID: req.UserID,
		PersonalDetails: entity.PersonalDetails{
			FullName:      req.PersonalDetails.FullName,
			PhoneNumber:   req.PersonalDetails.PhoneNumber,
			Email:         req.PersonalDetails.Email,
			Linkedin:      req.PersonalDetails.Linkedin,
			PortfolioURL:  req.PersonalDetails.PortfolioURL,
			Description:   req.PersonalDetails.Description,
			AddressString: req.PersonalDetails.AddressString,
		},
		ProfessionalExperience: []entity.Experience{},
		Education:              []entity.Education{},
		LeadershipExperience:   []entity.Leadership{},
		Others:                 []entity.Achievement{},
	}

	for _, exp := range req.ProfessionalExperience {
		result.ProfessionalExperience = append(result.ProfessionalExperience, entity.Experience{
			StartDate:   entity.Date{Month: exp.StartDate.Month, Year: exp.StartDate.Year},
			EndDate:     entity.Date{Month: exp.EndDate.Month, Year: exp.EndDate.Year},
			RoleTitle:   exp.RoleTitle,
			CompanyName: exp.CompanyName,
			Location:    exp.Location,
			Current:     exp.Current,
			Elaboration: []entity.Elaboration{},
		})

		for _, elab := range exp.Elaboration {
			result.ProfessionalExperience[len(result.ProfessionalExperience)-1].Elaboration = append(
				result.ProfessionalExperience[len(result.ProfessionalExperience)-1].Elaboration,
				entity.Elaboration{Text: elab.Text},
			)
		}
	}

	for _, edu := range req.Education {
		result.Education = append(result.Education, entity.Education{
			StartDate:   entity.Date{Month: edu.StartDate.Month, Year: edu.StartDate.Year},
			EndDate:     entity.Date{Month: edu.EndDate.Month, Year: edu.EndDate.Year},
			School:      edu.School,
			Location:    edu.Location,
			DegreeLevel: edu.DegreeLevel,
			Major:       edu.Major,
			GPA:         edu.GPA,
			MaxGPA:      edu.MaxGPA,
			Elaboration: []entity.Elaboration{},
		})

		for _, elab := range edu.Elaboration {
			result.Education[len(result.Education)-1].Elaboration = append(
				result.Education[len(result.Education)-1].Elaboration,
				entity.Elaboration{Text: elab.Text},
			)
		}
	}

	for _, lead := range req.LeadershipExperience {
		result.LeadershipExperience = append(result.LeadershipExperience, entity.Leadership{
			StartDate:        entity.Date{Month: lead.StartDate.Month, Year: lead.StartDate.Year},
			EndDate:          entity.Date{Month: lead.EndDate.Month, Year: lead.EndDate.Year},
			RoleTitle:        lead.RoleTitle,
			OrganisationName: lead.OrganisationName,
			Location:         lead.Location,
			Current:          lead.Current,
			Elaboration:      []entity.Elaboration{},
		})

		for _, elab := range lead.Elaboration {
			result.LeadershipExperience[len(result.LeadershipExperience)-1].Elaboration = append(
				result.LeadershipExperience[len(result.LeadershipExperience)-1].Elaboration,
				entity.Elaboration{Text: elab.Text},
			)
		}
	}

	for _, achievement := range req.Others {
		result.Others = append(result.Others, entity.Achievement{
			Elaboration: entity.Elaboration{Text: achievement.Elaboration.Text},
			Name:        achievement.Name,
			Date:        entity.Date{Month: achievement.Date.Month, Year: achievement.Date.Year},
			Category:    achievement.Category,
		})
	}

	return result, nil
}

func (h *ResumeHandler) isPDF(fileName string) bool {
	regex := regexp.MustCompile(`(?i)\.pdf$`)

	return regex.MatchString(fileName)
}
