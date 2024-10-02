package resumeService

import (
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/entity"
)

func (s resumeService) formattedResume(resumeData entity.Resume) resume.GetResumeResponse {
	return resume.GetResumeResponse{
		ID:        resumeData.ID,
		Name:      resumeData.Name,
		CreatedAt: resumeData.CreatedAt.String(),
	}
}

func (s resumeService) formattedResumeDetail(resumeDetail entity.ResumeDetail) resume.ResumeDetailDTO {
	// Convert PersonalDetails
	personalDetails := resume.PersonalDetails{
		FullName:      resumeDetail.PersonalDetails.FullName,
		PhoneNumber:   resumeDetail.PersonalDetails.PhoneNumber,
		Email:         resumeDetail.PersonalDetails.Email,
		Linkedin:      resumeDetail.PersonalDetails.Linkedin,
		PortfolioURL:  resumeDetail.PersonalDetails.PortfolioURL,
		Description:   resumeDetail.PersonalDetails.Description,
		AddressString: resumeDetail.PersonalDetails.AddressString,
	}

	// Convert ProfessionalExperience
	var professionalExperience []resume.Experience
	for _, exp := range resumeDetail.ProfessionalExperience {
		professionalExperience = append(professionalExperience, resume.Experience{
			StartDate:   resume.Date(exp.StartDate),
			EndDate:     resume.Date(exp.EndDate),
			RoleTitle:   exp.RoleTitle,
			CompanyName: exp.CompanyName,
			Location:    exp.Location,
			Current:     exp.Current,
			Elaboration: convertElaboration(exp.Elaboration),
		})
	}

	// Convert Education
	var education []resume.Education
	for _, edu := range resumeDetail.Education {
		education = append(education, resume.Education{
			StartDate:   resume.Date(edu.StartDate),
			EndDate:     resume.Date(edu.EndDate),
			School:      edu.School,
			Location:    edu.Location,
			DegreeLevel: edu.DegreeLevel,
			Major:       edu.Major,
			GPA:         edu.GPA,
			MaxGPA:      edu.MaxGPA,
			Elaboration: convertElaboration(edu.Elaboration),
		})
	}

	// Convert LeadershipExperience
	var leadershipExperience []resume.Leadership
	for _, lead := range resumeDetail.LeadershipExperience {
		leadershipExperience = append(leadershipExperience, resume.Leadership{
			StartDate:        resume.Date(lead.StartDate),
			EndDate:          resume.Date(lead.EndDate),
			RoleTitle:        lead.RoleTitle,
			OrganisationName: lead.OrganisationName,
			Location:         lead.Location,
			Current:          lead.Current,
			Elaboration:      convertElaboration(lead.Elaboration),
		})
	}

	// Convert Others (Achievements)
	var others []resume.Achievement
	for _, ach := range resumeDetail.Others {
		others = append(others, resume.Achievement{
			Name:        ach.Name,
			Date:        resume.Date(ach.Date),
			Category:    ach.Category,
			Elaboration: resume.Elaboration{Text: ach.Elaboration.Text},
		})
	}

	// Return the constructed ResumeDTO
	return resume.ResumeDetailDTO{
		ID:                     resumeDetail.ID.Hex(),
		UserID:                 resumeDetail.UserID,
		PersonalDetails:        personalDetails,
		ProfessionalExperience: professionalExperience,
		Education:              education,
		LeadershipExperience:   leadershipExperience,
		Others:                 others,
	}
}

// Helper function to convert Elaboration
func convertElaboration(elaborations []entity.Elaboration) []resume.Elaboration {
	var result []resume.Elaboration
	for _, el := range elaborations {
		result = append(result, resume.Elaboration{
			Text: el.Text,
		})
	}
	return result
}
