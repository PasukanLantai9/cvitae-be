package google

import (
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

func (g googleGemini) GenerateResumeJsonFromPDF(ctx context.Context, pdfFile []byte) (string, error) {
	content, err := g.model.GenerateContent(ctx,
		genai.Text(`I am providing a resume in the form of a PDF. Please extract the content and convert it into JSON with the following structure:
        {
            "personalDetails": {
                "fullName": "string",
                "phoneNumber": "string",
                "email": "string",
                "linkedin": "string",
                "portfolioUrl": "string",
                "description": "string",
                "addressString": "string"
            },
            "education": [
                {
                    "startDate": { "month": "string", "year": "int" },
                    "endDate": { "month": "string", "year": "int" },
                    "school": "string",
                    "location": "string",
                    "degreeLevel": "string",
                    "major": "string",
                    "gpa": 0.0,
                    "maxGpa": 0.0
                }
            ],
            "professionalExperience": [
                {
                    "startDate": { "month": "string", "year": "int" },
                    "endDate": { "month": "string", "year": "int" },
                    "roleTitle": "string",
                    "companyName": "string",
                    "location": "string",
                    "current": true
                }
            ],
            "leadershipExperience": [
                {
                    "startDate": { "month": "string", "year": "int" },
                    "endDate": { "month": "string", "year": "int" },
                    "roleTitle": "string",
                    "organisationName": "string",
                    "location": "string"
                }
            ],
            "others": [
                {
                    "name": "string",
                    "date": { "month": "string", "year": "int" },
                    "category": "string",
                    "elaboration": { "text": "string" }
                }
            ]
        }`),
		genai.Text("Make sure all fields are correctly formatted, especially dates and categories like 'education', 'professionalExperience', etc."),
		genai.Text("Make sure the 'others.category' field only contains 'Certificates', 'Skills', or 'Achievements'."),
		genai.Text("If the PDF content does not match a resume format, return 'not-resume' as the output instead."),
		genai.Text("Give me only the JSON output in one-line, without anything else"),
		genai.Blob{
			MIMEType: "application/pdf",
			Data:     pdfFile,
		},
	)
	if err != nil {
		return "", err
	}

	part := content.Candidates[0].Content.Parts[0]
	byteJson, err := json.Marshal(part)
	if err != nil {
		return "", err
	}

	strJson, err := strconv.Unquote(string(byteJson))
	if err != nil {
		return "", err
	}

	return strJson, nil
}

func (g googleGemini) GenerateExperienceAndSkillsParagrafFromJSON(ctx context.Context, jsonData []byte) (string, error) {
	content, err := g.model.GenerateContent(ctx,
		genai.Text(`Based on the following JSON data, generate a CV summary paragraph that combines the user's professional experience and skills. For example: 
"I have experience as a Graphic Designer for over 5 years, working at PT. Kreatif Abadi in Jakarta. During my time there, I developed marketing materials such as brochures, posters, and social media content, collaborating with the marketing team to create engaging campaigns. 
My skills include proficiency in design software like Adobe Photoshop, Illustrator, and InDesign, as well as experience in UI/UX design, project management, and digital illustration."`),
		genai.Text("Ensure all relevant data from the JSON is included, and the output should be in one line."),
		genai.Text("If the PDF content does not match a resume format, return 'not-resume' as the output instead."),
		genai.Text("Give me only the string output in one line, without any additional formatting or extraneous characters."),
		genai.Text(jsonData),
	)
	if err != nil {
		return "", err
	}

	var formatted strings.Builder
	if content != nil && content.Candidates != nil {
		for _, cand := range content.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formatted.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	return formatted.String(), nil
}

func (g googleGemini) GenerateExperienceAndSkillsParagrafFromPDF(ctx context.Context, pdfFile []byte) (string, error) {
	content, err := g.model.GenerateContent(ctx,
		genai.Text(`Based on the following CV PDF data, generate a CV summary paragraph that combines the user's professional experience and skills. For example: 
"I have experience as a Graphic Designer for over 5 years, working at PT. Kreatif Abadi in Jakarta. During my time there, I developed marketing materials such as brochures, posters, and social media content, collaborating with the marketing team to create engaging campaigns. 
My skills include proficiency in design software like Adobe Photoshop, Illustrator, and InDesign, as well as experience in UI/UX design, project management, and digital illustration."`),
		genai.Text("Ensure all relevant data from the CV PDF is included, and the output should be in one line."),
		genai.Text("If the PDF content does not match a resume format, return 'not-resume' as the output instead."),
		genai.Text("Give me only the string output in one line, without any additional formatting or extraneous characters."),
		genai.Blob{
			MIMEType: "application/pdf",
			Data:     pdfFile,
		},
	)
	if err != nil {
		return "", err
	}

	var formatted strings.Builder
	if content != nil && content.Candidates != nil {
		for _, cand := range content.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formatted.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	return formatted.String(), nil
}
