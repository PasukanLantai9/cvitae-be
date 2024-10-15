package google

import (
	"encoding/json"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/context"
	"strconv"
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
