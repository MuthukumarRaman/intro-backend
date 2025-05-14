package openaiservice

import (
	"context"
	"encoding/json"
	"fmt"

	"introme-api/pkg/shared/database"
	"introme-api/pkg/shared/helper"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var openAiKey = "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"
var openAiOrgId = "org-CmUrsek5G1rJm0RYVMX6om1B"

// AIProfile represents an individual profile configuration
type AIProfile struct {
	ModelID        string       `json:"modelId"`
	TaskDefinition string       `json:"taskDefinition"`
	AIFunctions    []AIFunction `json:"aiFunctions"`
	ForcedFunction string       `json:"forcedFunction"`
}

// AIConfigModel represents the overall configuration for AI profiles
type AIConfigModel struct {
	Profiles map[string]AIConfigProfile `json:"profiles"`
}

// AIConfigProfile represents an individual profile configuration
type AIConfigProfile struct {
	ModelID        string       `json:"modelId"`
	TaskDefinition string       `json:"taskDefinition"`
	AIFunctions    []AIFunction `json:"aiFunctions"`
	ForcedFunction string       `json:"forcedFunction"`
}

// AIFunction represents a function configuration for AI profiles
type AIFunction struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Parameters  Property `json:"parameters"`
}

type OpenAIDescriptors struct{}

// Property defines a single property with type and description
type Property struct {
	Type        string              `json:"type"`
	Format      string              `json:"format"`
	Category    string              `json:"category,omitempty"`
	Description string              `json:"description,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Default     string              `json:"default,omitempty"`
}

func (d *OpenAIDescriptors) OpenAIDescriptorsConfig() Property {

	// return Property{
	// 	Type: "object",
	// 	Properties: map[string]Property{
	// 		"userProfile": {
	// 			Type:        "object",
	// 			Description: "A comprehensive profile where you describe yourself.",
	// 			Required:    []string{"first_name", "email_id", "gender", "profile_image", "last_name", "my_intro", "bio", "age", "location", "industry", "work_life_philosophy", "professional_journey", "expertise", "hobbies"},
	// 			Properties: map[string]Property{
	// 				"my_intro": {
	// 					Type:        "string",
	// 					Description: "Write a short introduction about yourself.",
	// 					Default:     "I am an enthusiastic professional eager to learn and grow.",
	// 				},
	// 				"bio": {
	// 					Type:        "string",
	// 					Description: "Describe yourself briefly, including your profession and experience.",
	// 					Default:     "I am a dedicated professional with a strong background in my field.",
	// 				},
	// 				"age": {
	// 					Type:        "number",
	// 					Description: "Enter your age, calculated from your date of birth.",
	// 					Default:     "25",
	// 				},
	// 				"location": {
	// 					Type:        "string",
	// 					Description: "Mention your current location (city and country).",
	// 					Default:     "Not specified",
	// 				},
	// 				"industry": {
	// 					Type:        "string",
	// 					Description: "List the industries you have experience in.",
	// 					Default:     "Technology",
	// 				},
	// 				"work_life_philosophy": {
	// 					Type:        "string",
	// 					Description: "Share your thoughts on work-life balance and your professional approach.",
	// 					Default:     "I believe in maintaining a healthy balance between work and personal life.",
	// 				},
	// 				"professional_journey": {
	// 					Type:        "string",
	// 					Description: "Summarize your career path and what led you to your current role.",
	// 					Default:     "I started my career with a passion for problem-solving and grew into my current role through continuous learning.",
	// 				},
	// 				"expertise": {
	// 					Type:        "string",
	// 					Description: "List your key skills and areas of expertise.",
	// 					Default:     "Problem-solving, Communication, Technical Skills",
	// 				},
	// 				"hobbies": {
	// 					Type:        "string",
	// 					Description: "Mention your hobbies and activities outside of work.",
	// 					Default:     "Reading, Traveling, Fitness",
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	return Property{
		Type: "object",
		Properties: map[string]Property{
			"userProfile": {
				Type:        "object",
				Description: "A comprehensive profile where you describe yourself.",
				Required:    []string{"first_name", "key_skills", "professional_match_criteria", "email_id", "birth_date", "gender", "profile_image", "experience_years", "last_name", "my_intro", "bio", "age", "location", "industry", "work_life_philosophy", "professional_journey", "expertise", "hobbies"},
				Properties: map[string]Property{
					"first_name": {
						Type:        "string",
						Description: "Your given name.",
					},
					"last_name": {
						Type:        "string",
						Description: "Your family name or surname.",
					},
					"email_id": {
						Type:        "string",
						Description: "Your primary email address.",
					},
					"gender": {
						Type:        "string",
						Description: "Your gender identity.",
					},

					"key_skills": {
						Type:        "array",
						Description: "Your Skills identity.",
						Items: &Property{
							Type: "string",
						},
					},
					"birth_date": {
						Type:        "string",
						Description: "Your date of birth.",
					},
					"professional_match_criteria": {
						Type:        "string",
						Description: "User's professional match criteria",
					},
					"profile_image": {
						Type:        "string",
						Description: "A URL linking to your profile picture.",
					},
					"my_intro": {
						Type:        "string",
						Description: "Write a short introduction about yourself. Share a brief overview of who you are, your passions, and what drives you. Mention your core values and what makes you unique. Highlight your enthusiasm for learning, growing, and contributing to your field. Provide insight into how you approach challenges and what excites you the most in your journey.",
					},
					"bio": {
						Type:        "string",
						Description: "Describe yourself briefly, including your profession and experience. Provide an overview of your background, the industries you have worked in, and the skills you have acquired. Highlight key achievements or milestones that have shaped your professional journey. Discuss your strengths and the areas where you excel. Share how your experience and expertise contribute to your current role and aspirations for the future.",
					},
					"age": {
						Type:        "number",
						Description: "Enter your age, calculated from your date of birth.",
					},
					"experience_years": {
						Type:        "number",
						Description: "experience of the user",
					},
					"location": {
						Type:        "string",
						Description: "Automatically filled with the location name based on the geo_location.",
					},
					"geo_location": {
						Type:        "array",
						Description: "Automatically filled with the location coordinates (longitude, latitude).",
						Items: &Property{
							Type: "number",
						},
					},
					"industries": {
						Type:        "array",
						Description: "Fill the industry Worked",
						Items: &Property{
							Type: "string",
						},
					},

					"work_life_philosophy": {
						Type:        "string",
						Description: "Share your thoughts on work-life balance and your professional approach.",
					},
					"professional_journey": {
						Type:        "string",
						Description: "Summarize your career path and what led you to your current role.elabrate it",
					},
					"expertise": {
						Type:        "string",
						Description: "List your key skills and areas of expertise.elabrate it",
					},
					"hobbies": {
						Type:        "string",
						Description: "Mention your hobbies and activities outside of work.",
					},
					"embedding_keywords": {
						Type:        "string",
						Description: "In this extract the important keys in this profile and stored it in this field",
					},
				},
			},
			"match_profile": {
				Type:        "object",
				Description: "Users Data",
				Properties: map[string]Property{
					"users": {
						Type: "array",
						Items: &Property{
							Type: "object",
							Properties: map[string]Property{
								"id": {
									Type:        "string",
									Description: "Id of the user",
								},

								"interests": {
									Type: "array",
									Items: &Property{
										Type: "string",
									},
									Description: "Interests of the user",
								},
							},
						},
					},
				},
			},
			"onboardingQuestions": {
				Type:        "object",
				Description: "Return an onboarding guide as natural language using a structured JSON model. Each question should have a `question` field for the prompt text and a `type` field for the input type (e.g., input, date, textarea, radio). All done by category wise",
				Properties: map[string]Property{
					"questions": {
						Type:        "array",
						Description: "A list of onboarding questions with their details in category wise.",
						Items: &Property{
							Type:        "object",
							Description: "A natural language question to ask the user  details in category wise",
							Properties: map[string]Property{
								"question": {
									Type:        "string",
									Description: "The text of the question to ask the user.",
									Category:    "Personal Details",
								},
								"type": {
									Type:        "string",
									Description: "The type of input for the question (e.g., input, date, textarea, radio).",
								},
								"sampleAnswer": {
									Type:        "string",
									Description: "An example answer for the question.",
								},
								"order": {
									Type:        "number",
									Description: "Order Vise view of question",
								},
								"category": {
									Type:        "string",
									Description: "questions category like educational details etc",
								},
								"field_name": {
									Type:        "string",
									Description: "Generated Question must be to concurrent field name like first_name etc",
								},
								"dataType": {
									Type: "string",
									// enum: ["string", "number", "date", "boolean"],
									Description: "The type of the answer (string, number, date,textarea, radio,chip or etc).",
								},
								"questions": {
									Type:        "object",
									Description: "A comprehensive profile where you describe yourself.",
									Required:    []string{"first_name", "key_skills", "professional_match_criteria", "email_id", "birth_date", "gender", "profile_image", "experience_years", "last_name", "my_intro", "bio", "age", "location", "industry", "work_life_philosophy", "professional_journey", "expertise", "hobbies"},
									Properties: map[string]Property{
										"first_name": {
											Type:        "string",
											Description: "Your given name.",
											Category:    "Personal Details",
										},
										"last_name": {
											Type:        "string",
											Description: "Your family name or surname.",
											Category:    "Personal Details",
										},
										"middle_name": {
											Type:        "string",
											Description: "Your family name or surname.",
											Category:    "Personal Details",
										},
										"birth_date": {
											Type:        "string",
											Description: "Your date of birth.",
											Category:    "Personal Details",
										},
										"email_id": {
											Type:        "string",
											Description: "Your primary email address.",
											Category:    "Personal Details",
										},
										"gender": {
											Type:        "string",
											Description: "Your gender identity.",
											Category:    "Personal Details",
										},
										"age": {
											Type:        "number",
											Description: "Enter your age, calculated from your date of birth.",
											Category:    "Personal Details",
										},
										"key_skills": {
											Type:        "array",
											Description: "Your Skills identity.",
											Items: &Property{
												Type: "string",
											},
											Category: "Educational Details",
										},

										"professional_match_criteria": {
											Type:        "string",
											Description: "User's professional match criteria",
											Category:    "Educational Details",
										},
										"profile_image": {
											Type:        "string",
											Description: "A URL linking to your profile picture.",
										},
										"my_intro": {
											Type:        "string",
											Description: "Write a short introduction about yourself. Share a brief overview of who you are, your passions, and what drives you. Mention your core values and what makes you unique. Highlight your enthusiasm for learning, growing, and contributing to your field. Provide insight into how you approach challenges and what excites you the most in your journey.",
										},
										"bio": {
											Type:        "string",
											Description: "Describe yourself briefly, including your profession and experience. Provide an overview of your background, the industries you have worked in, and the skills you have acquired. Highlight key achievements or milestones that have shaped your professional journey. Discuss your strengths and the areas where you excel. Share how your experience and expertise contribute to your current role and aspirations for the future.",
										},

										"experience_years": {
											Type:        "number",
											Description: "experience of the user",
										},
										"location": {
											Type:        "string",
											Description: "Automatically filled with the location name based on the geo_location.",
										},
										"geo_location": {
											Type:        "array",
											Description: "Automatically filled with the location coordinates (longitude, latitude).",
											Items: &Property{
												Type: "number",
											},
										},
										"industries": {
											Type:        "array",
											Description: "Fill the industry Worked",
											Items: &Property{
												Type: "string",
											},
										},
										"work_life_philosophy": {
											Type:        "string",
											Description: "Share your thoughts on work-life balance and your professional approach.",
										},
										"professional_journey": {
											Type:        "string",
											Description: "Summarize your career path and what led you to your current role.elabrate it",
										},
										"expertise": {
											Type:        "string",
											Description: "List your key skills and areas of expertise.elabrate it",
										},
										"hobbies": {
											Type:        "string",
											Description: "Mention your hobbies and activities outside of work.",
										},
									},
								},
							},
						},
					},
				},
			},
			// "onboardingQuestionsWithModule": {
			// 	Type:        "object",
			// 	Description: "Return an onboarding guide as natural language using a structured JSON model. Each question should have a `question` field for the prompt text and a `type` field for the input type (e.g., input, date, textarea, radio).",
			// 	Properties: map[string]Property{
			// 		Items: &Property{
			// 			"Module": {
			// 				Type:        "object",
			// 				Description: "A list of onboarding questions with their details.",
			// 				Properties: map[string]Property{
			// 					"personalDetails": {
			// 						Type:        "object",
			// 						Description: "A list of onboarding questions with their details.",
			// 						Items: &Property{
			// 							Type:        "object",
			// 							Description: "A natural language question to ask the user",
			// 							Properties: map[string]Property{
			// 								"question": {
			// 									Type:        "string",
			// 									Description: "The text of the question to ask the user.",
			// 								},
			// 								"type": {
			// 									Type:        "string",
			// 									Description: "The type of input for the question (e.g., input, date, textarea, radio).",
			// 								},
			// 								"sampleAnswer": {
			// 									Type:        "string",
			// 									Description: "An example answer for the question.",
			// 								},
			// 								"order": {
			// 									Type:        "number",
			// 									Description: "Order Vise view of question",
			// 								},
			// 								"dataType": {
			// 									Type: "string",
			// 									// enum: ["string", "number", "date", "boolean"],
			// 									Description: "The type of the answer (string, number, date,textarea, radio,chip or etc).",
			// 								},
			// 							},
			// 						},
			// 					},
			// 					// "personalDetails": [
			// 					//   {
			// 					// 	"question": "What is your full name?",
			// 					// 	"type": "input",
			// 					// 	"sampleAnswer": "John Doe",
			// 					// 	"order": 1,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "What is your date of birth?",
			// 					// 	"type": "date",
			// 					// 	"sampleAnswer": "1990-01-01",
			// 					// 	"order": 2,
			// 					// 	"dataType": "date"
			// 					//   },
			// 					//   {
			// 					// 	"question": "What is your contact number?",
			// 					// 	"type": "input",
			// 					// 	"sampleAnswer": "+1234567890",
			// 					// 	"order": 3,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "What is your email address?",
			// 					// 	"type": "input",
			// 					// 	"sampleAnswer": "john.doe@example.com",
			// 					// 	"order": 4,
			// 					// 	"dataType": "string"
			// 					//   }
			// 					// ],
			// 					// "educationalDetails": [
			// 					//   {
			// 					// 	"question": "What is your highest qualification?",
			// 					// 	"type": "input",
			// 					// 	"sampleAnswer": "Master’s in Computer Science",
			// 					// 	"order": 1,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "Which university/college did you attend?",
			// 					// 	"type": "input",
			// 					// 	"sampleAnswer": "Harvard University",
			// 					// 	"order": 2,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "What was your field of study?",
			// 					// 	"type": "input",
			// 					// 	"sampleAnswer": "Computer Science",
			// 					// 	"order": 3,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "What was your year of graduation?",
			// 					// 	"type": "date",
			// 					// 	"sampleAnswer": "2015",
			// 					// 	"order": 4,
			// 					// 	"dataType": "date"
			// 					//   }
			// 					// ],
			// 					// "personalInterests": [
			// 					//   {
			// 					// 	"question": "What are your favorite hobbies?",
			// 					// 	"type": "chip",
			// 					// 	"sampleAnswer": "Reading, Traveling, Coding",
			// 					// 	"order": 1,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "Do you enjoy outdoor activities?",
			// 					// 	"type": "radio",
			// 					// 	"sampleAnswer": "Yes",
			// 					// 	"order": 2,
			// 					// 	"dataType": "boolean"
			// 					//   },
			// 					//   {
			// 					// 	"question": "What type of books do you like to read?",
			// 					// 	"type": "textarea",
			// 					// 	"sampleAnswer": "Science fiction, Mystery, Self-improvement",
			// 					// 	"order": 3,
			// 					// 	"dataType": "string"
			// 					//   },
			// 					//   {
			// 					// 	"question": "Are you interested in learning new skills?",
			// 					// 	"type": "radio",
			// 					// 	"sampleAnswer": "Yes",
			// 					// 	"order": 4,
			// 					// 	"dataType": "boolean"
			// 					//   }
			// 					// ]
			// 				},
			// 			},
			// 		}},
			// },
			// 		"userProfile": {
			// 			Type:        "object",
			// 			Description: "A comprehensive profile where you describe yourself.",
			// 			Required:    []string{"first_name", "email_id", "gender", "profile_image", "last_name", "my_intro", "bio", "age", "location", "industry", "work_life_philosophy", "professional_journey", "expertise", "hobbies"},
			// 			Properties: map[string]Property{
			// 				"my_intro": {
			// 					Type:        "string",
			// 					Description: "Write a short introduction about yourself.",
			// 					Default:     "I am an enthusiastic professional eager to learn and grow.",
			// 				},
			// 				"bio": {
			// 					Type:        "string",
			// 					Description: "Describe yourself briefly, including your profession and experience.",
			// 					Default:     "I am a dedicated professional with a strong background in my field.",
			// 				},
			// 				"age": {
			// 					Type:        "number",
			// 					Description: "Enter your age, calculated from your date of birth.",
			// 					Default:     "25",
			// 				},
			// 				"location": {
			// 					Type:        "string",
			// 					Description: "Mention your current location (city and country).",
			// 					Default:     "Not specified",
			// 				},
			// 				"industry": {
			// 					Type:        "string",
			// 					Description: "List the industries you have experience in.",
			// 					Default:     "Technology",
			// 				},
			// 				"work_life_philosophy": {
			// 					Type:        "string",
			// 					Description: "Share your thoughts on work-life balance and your professional approach.",
			// 					Default:     "I believe in maintaining a healthy balance between work and personal life.",
			// 				},
			// 				"professional_journey": {
			// 					Type:        "string",
			// 					Description: "Summarize your career path and what led you to your current role.",
			// 					Default:     "I started my career with a passion for problem-solving and grew into my current role through continuous learning.",
			// 				},
			// 				"expertise": {
			// 					Type:        "string",
			// 					Description: "List your key skills and areas of expertise.",
			// 					Default:     "Problem-solving, Communication, Technical Skills",
			// 				},
			// 				"hobbies": {
			// 					Type:        "string",
			// 					Description: "Mention your hobbies and activities outside of work.",
			// 					Default:     "Reading, Traveling, Fitness",
			// 				},
			// 			},
			// 		},
		},
	}

}

// NewAIConfigModel initializes the AI configuration model
func NewAIConfigModel(descriptors *OpenAIDescriptors) AIConfigModel {
	return AIConfigModel{
		Profiles: map[string]AIConfigProfile{
			// "onboarding": {
			// 	ModelID:        "gpt-3.5-turbo",
			// 	TaskDefinition: "Return an onboarding guide as natural language using a JSON structure.",
			// 	AIFunctions: []AIFunction{
			// 		{
			// 			Name:        "parseToOnboardingModel",
			// 			Description: "Parse an onboarding guide description to a JSON model",
			// 			Parameters:  descriptors.OpenAIDescriptorsConfig(),
			// 		},
			// 	},
			// 	ForcedFunction: "parseToOnboardingModel",
			// },
			// "onboardingQuestions": {
			// 	ModelID:        "gpt-3.5-turbo",
			// 	TaskDefinition: "Return an onboarding guide as natural language using a structured JSON model. Each question should have a `question` field for the prompt text and a `type` field for the input type (e.g., input, date, textarea, radio).",
			// 	AIFunctions: []AIFunction{
			// 		{
			// 			Name:        "parseToOnboardingModel",
			// 			Description: "Parse an onboarding guide description to a JSON model",
			// 			Parameters:  descriptors.OpenAIDescriptorsConfig(),
			// 		},
			// 	},
			// 	ForcedFunction: "parseToOnboardingModel",
			// },
			"onboardingQuestions": {
				ModelID:        "gpt-3.5-turbo",
				TaskDefinition: "Return an onboarding guide as natural language using a structured JSON model. Each question should have a `question` field for the prompt text and a `type` field for the input type (e.g., input, date, textarea, radio,number).",
				AIFunctions: []AIFunction{
					{
						Name:        "parseToOnboardingModel",
						Description: "Parse an onboarding guide description to a JSON model",
						Parameters:  descriptors.OpenAIDescriptorsConfig().Properties["onboardingQuestions"],
					},
				},
				ForcedFunction: "parseToOnboardingModel",
			},
			"userProfile": {
				ModelID:        "gpt-4o", //gpt-3.5-turbo
				TaskDefinition: "Return a user profile as natural language using a JSON structure.",
				AIFunctions: []AIFunction{
					{
						Name:        "parseToUserProfileModel",
						Description: "Parse a user profile description to a JSON model",
						Parameters:  descriptors.OpenAIDescriptorsConfig().Properties["userProfile"],
					},
				},
				ForcedFunction: "parseToUserProfileModel",
			},
			"match_profile": {
				ModelID:        "gpt-3.5-turbo", //gpt-4o
				TaskDefinition: "Return only the matched users by strictly comparing their interests  using a structured JSON format. DO NOT RETURN USERS WHO DO NOT MATCH.",
				AIFunctions: []AIFunction{
					{
						Name:        "matchUserProfileModel",
						Description: "Filter and return ONLY the users whose interests match .  Return the user IDs only.",
						Parameters:  descriptors.OpenAIDescriptorsConfig().Properties["match_profile"],
					},
				},
				ForcedFunction: "matchUserProfileModel",
			},
			"match_reason": {
				ModelID:        "gpt-3.5-turbo", //gpt-4o
				TaskDefinition: "The match_reason must be exactly three lines long, providing a clear explanation of why the match was made.",
				AIFunctions: []AIFunction{
					{
						Name:        "matchProfileReasonModel",
						Description: "The match_reason must be exactly three lines long, providing a clear explanation of why the match was made.",
						Parameters:  descriptors.OpenAIDescriptorsConfig().Properties["match_profile"],
					},
				},
				ForcedFunction: "matchProfileReasonModel",
			},
		},
	}
}

func GenerateFromAI(client *openai.Client, aiQuery string, targetConfig string, descriptors *OpenAIDescriptors) (map[string]interface{}, error) {
	// log.Println("[OPENAI] Calling OpenAI service")

	config := NewAIConfigModel(descriptors)
	aiConfig, exists := config.Profiles[targetConfig]
	if !exists {
		return nil, fmt.Errorf("targetConfig '%s' not found", targetConfig)
	}

	// Prepare the query and OpenAI request
	query := aiConfig.TaskDefinition + "\n" + aiQuery
	modelID := aiConfig.ModelID
	aiFunctions := aiConfig.AIFunctions
	forcedFunction := aiConfig.ForcedFunction

	// log.Printf("[OPENAI] Query: %s\n", query)
	// Convert your AIFunction slice to openai.FunctionDefinition slice
	var openAIFunctions []openai.FunctionDefinition
	for _, f := range aiFunctions {
		openAIFunctions = append(openAIFunctions, openai.FunctionDefinition{
			Name:        f.Name,
			Description: f.Description,
			Parameters:  f.Parameters,
		})
	}

	// Create OpenAI Chat request
	req := openai.ChatCompletionRequest{
		Model: modelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: query},
		},
		Functions:    openAIFunctions,
		FunctionCall: openai.FunctionCall{Name: forcedFunction},
	}

	ctx := context.Background()
	// Perform the OpenAI API call
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Parse and return the function call arguments
	var result map[string]interface{}
	if len(resp.Choices) > 0 && resp.Choices[0].Message.FunctionCall != nil {
		err = json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
func GenerateEmbeddingFromAI(client *openai.Client, text string) ([]float32, error) {
	// Prepare the OpenAI Embedding request
	ctx := context.Background()
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: []string{text},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %v", err)
	}

	// Extract embedding vector
	if len(resp.Data) > 0 {

		return resp.Data[0].Embedding, nil
	}

	return nil, fmt.Errorf("no embedding returned")
}

func ProfileMatchFromAI(client *openai.Client, aiQuery string, targetConfig string, descriptors *OpenAIDescriptors) ([]string, error) {
	// log.Println("[OPENAI] Calling OpenAI service")

	config := NewAIConfigModel(descriptors)
	aiConfig, exists := config.Profiles[targetConfig]
	if !exists {
		return nil, fmt.Errorf("targetConfig '%s' not found", targetConfig)
	}

	// Prepare the query and OpenAI request
	query := aiConfig.TaskDefinition + "\n" + aiQuery
	modelID := aiConfig.ModelID
	aiFunctions := aiConfig.AIFunctions
	forcedFunction := aiConfig.ForcedFunction

	// Convert your AIFunction slice to openai.FunctionDefinition slice
	var openAIFunctions []openai.FunctionDefinition
	for _, f := range aiFunctions {
		openAIFunctions = append(openAIFunctions, openai.FunctionDefinition{
			Name:        f.Name,
			Description: f.Description,
			Parameters:  f.Parameters,
		})
	}

	// Create OpenAI Chat request
	req := openai.ChatCompletionRequest{
		Model: modelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: query},
		},
		Functions:    openAIFunctions,
		FunctionCall: openai.FunctionCall{Name: forcedFunction},
	}

	ctx := context.Background()
	// Perform the OpenAI API call
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Extract and return user IDs with matching skills, hobbies, and expertise
	var result struct {
		Matches []struct {
			UserID string `json:"user_id"`
		} `json:"matches"`
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message.FunctionCall != nil {
		err = json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &result)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(resp.Choices[0].Message.FunctionCall)
	// Collect user IDs
	var userIDs []string
	for _, match := range result.Matches {
		userIDs = append(userIDs, match.UserID)
	}

	return userIDs, nil
}

func ProfileMatchFromOpenAI(client *openai.Client, aiQuery string, targetConfig string, descriptors *OpenAIDescriptors) (map[string]interface{}, error) {
	// log.Println("[OPENAI] Calling OpenAI service")

	config := NewAIConfigModel(descriptors)
	aiConfig, exists := config.Profiles[targetConfig]
	if !exists {
		return nil, fmt.Errorf("targetConfig '%s' not found", targetConfig)
	}

	// Prepare the query and OpenAI request
	// query := aiQuery
	modelID := aiConfig.ModelID
	// Create OpenAI Chat request
	req := openai.ChatCompletionRequest{
		Model: modelID,

		// Messages: []openai.ChatCompletionMessage{
		// 	{
		// 		Role:    "system",
		// 		Content: "You are an assistant that compares user key_skills. Return users from 'otherUsers' who have at least one matching key_skills with the 'primaryUser'.",
		// 	},
		// 	// {Role: "user", Content: query},
		// 	{
		// 		Role:    "user",
		// 		Content: aiQuery,
		// 	},
		// },
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: "You are an assistant that compares users. Return users from 'otherUsers' who have at least one matching attribute with the 'primaryUser'. " +
					"For each matched user, return a match_type and a match_reason. " +
					"The match_reason must be exactly three lines long, providing a clear explanation of why the match was made.",
			},
			{
				Role:    "user",
				Content: aiQuery,
			},
		},

		Tools: []openai.Tool{
			{
				Type: "function",
				Function: &openai.FunctionDefinition{
					Name:        "get_matched_users",
					Description: "Get Matched Users with match type and reason",
					Parameters: Property{
						Type: "object",
						Properties: map[string]Property{
							"matched_users": Property{
								Type: "array",
								Items: &Property{
									Type: "object",
									Properties: map[string]Property{
										"user_id": {
											Type:        "string",
											Description: "Matched user ID",
										},
										"match_type": {
											Type:        "string",
											Description: "Type of match (e.g., Synergy, Connection, Fit)",
										},
										"match_reason": {
											Type:        "string",
											Description: "Reason why this user is a match",
										},
									},
									Required: []string{"user_id", "match_type", "match_reason"},
								},
								Description: "Array of matched users with details",
							},
						},
						Required: []string{"matched_users"},
					},
				},
			},
		},
	}

	ctx := context.Background()
	// Perform the OpenAI API call
	fmt.Println("ai called......")
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Println("ai responsed.....")

	// Extract and return user IDs with matching skills, hobbies, and expertise

	var args map[string]interface{}

	if len(resp.Choices) > 0 &&
		len(resp.Choices[0].Message.ToolCalls) > 0 &&
		resp.Choices[0].Message.ToolCalls[0].Function.Arguments != "" {

		err := json.Unmarshal([]byte(resp.Choices[0].Message.ToolCalls[0].Function.Arguments), &args)
		if err != nil {
			return nil, err
		}
		return args, nil
	}
	if len(resp.Choices) == 0 &&
		len(resp.Choices[0].Message.ToolCalls) == 0 {

		return nil, nil
	}

	return nil, fmt.Errorf("Error Getting Response From GenAI")
}

type ResponseModel struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Details string      `json:"details"`
	Data    interface{} `json:"data,omitempty"`
}

func GetUserProfile2(c *fiber.Ctx) error {

	// Extract query parameter and body
	profileID := c.Params("profileId")
	var result map[string]interface{}
	var newProfileData map[string]interface{}
	if err := c.BodyParser(&newProfileData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseModel{
			Status:  "failure",
			Message: "Error parsing request body",
			Details: err.Error(),
		})
	}
	var newData = true
	// Call the controller
	if profileID != "" {

		DataUpdateById(newProfileData, profileID)
		newData = false
	}

	client := openai.NewClient(openAiKey)
	var strcut = OpenAIDescriptors{}
	var aiQuery string

	if newData {
		// Construct the query
		aiQuery = fmt.Sprintf(
			"Turn this user profile into a natural language description based on the following user profile model:\n%s\n"+
				"ONLY INTERPRET THE DATA IN THE USER PROFILE DATA PROVIDED. DO NOT MAKE UP ANYTHING ELSE.\n"+
				"The new information to add to the user profile is:\n%s\n",
			strcut.OpenAIDescriptorsConfig().Properties["userProfile"], newProfileData,
		)
	} else {
		// Construct the AI query string using fmt.Sprintf
		aiQuery = fmt.Sprintf(
			"Turn this user profile into a natural language description basd on the following user profile model: \n"+
				"%s"+
				"ONLY INTERPRET THE DATA IN THE USER PROFILE DATA PROVIDED. DO NOT MAKE UP ANYTHING ELSE. \n"+
				"The current user profile data is: \n"+"%s"+"\n"+
				"The new information to add to the user profile is: \n"+"%s"+"\n"+
				"",
			// Assuming OpenAIDescriptorsConfig().Properties["userProfile"] is serialized to a string
			strcut.OpenAIDescriptorsConfig().Properties["userProfile"], result, newProfileData,
		)
	}

	if profileID == "" {
		res, err := database.GetConnection().Collection("user").InsertOne(context.Background(), newProfileData)
		if err != nil {
			return helper.BadRequest("Failed to insert data into the database: " + err.Error())
		}
		profileID = res.InsertedID.(string)
	}

	res, err := GenerateFromAI(client, aiQuery, "userProfile", &strcut)
	if err != nil {
		// return shared.InternalServerError(err.Error())
	}

	DataUpdateById(res["userProfile"].(map[string]interface{}), profileID)

	return helper.SuccessResponse(c, res)
}

func MatchUserProfile(c *fiber.Ctx) error {

	var primaryUser = map[string]interface{}{
		"id":        420,
		"name":      "Siva",
		"email":     "siva@gmail.com",
		"age":       29,
		"location":  "New York",
		"interests": []string{"reading", "riding", "coding"},
	}

	// var users []map[string]interface{}
	var OtherUsers = []map[string]interface{}{
		{
			"id":        1,
			"name":      "John Doe",
			"email":     "john.doe@example.com",
			"age":       29,
			"location":  "New York",
			"hobbies":   []string{"gaming", "gambling"},
			"interests": []string{"reading", "hiking", "coding"},
		},
		{
			"id":        23,
			"name":      "Jane Smith",
			"email":     "jane.smith@example.com",
			"age":       34,
			"location":  "Los Angeles",
			"hobbies":   []string{"gaming", "gambling"},
			"interests": []string{"reading", "hiking", "coding"},
		},
		{
			"id":        3,
			"name":      "Alice Johnson",
			"email":     "alice.johnson@example.com",
			"age":       25,
			"location":  "Chicago",
			"interests": []string{"music", "yoga", "gaming"},
		},
		{
			"id":        4,
			"name":      "Alice Johnson",
			"email":     "alice.johnson@example.com",
			"age":       25,
			"location":  "Chicago",
			"interests": []string{"reading", "coding"},
		},
	}

	var strcut = OpenAIDescriptors{}

	// aiQuery := fmt.Sprintf(
	// 	"Find and return ONLY the users from the following list whose skills, hobbies, location, interests, or expertise match the specified criteria. "+
	// 		"STRICTLY EXCLUDE USERS WHO DO NOT MATCH WITH PRIMARY USER SKILLS,HOBBIES ETC.\n\n"+
	// 		"Primary user profile:\n%v\n\n"+
	// 		"List of other users to check for matches:\n%v\n\n"+
	// 		"Return only the user IDs in a structured JSON format.",
	// 	primaryUser, OtherUsers,
	// )

	aiQuery := fmt.Sprintf(
		"You are a strict filtering engine. Your task is to find and return ONLY the users from the following list whose interests have at least TWO matches with the primary user’s interests. "+
			"DO NOT return users with fewer than TWO common interests. "+
			"If no users match this condition, return an empty array. "+
			"STRICTLY FOLLOW THIS RULE.\n\n"+
			"Primary user profile:\n%v\n\n"+
			"List of other users to check for matches:\n%v\n\n"+
			"Output format (no extra text): ",
		primaryUser, OtherUsers,
	)

	client := openai.NewClient(openAiKey)
	res, err := ProfileMatchFromOpenAI(client, aiQuery, "match_profile", &strcut)
	if err != nil {
		return helper.InternalServerError(err.Error())
	}

	return helper.SuccessResponse(c, res)
}

func GetUserProfile(c *fiber.Ctx) error {
	profileID := c.Params("profileId")
	var newProfileData map[string]interface{}

	if err := c.BodyParser(&newProfileData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseModel{
			Status:  "failure",
			Message: "Error parsing request body",
			Details: err.Error(),
		})
	}

	newData := profileID == ""

	existingProfile1, _ := fetchUserProfile(profileID)
	if existingProfile1 == nil {
		newData = true
	}

	client := openai.NewClient(openAiKey)

	var descriptor OpenAIDescriptors
	aiQuery := buildAIQuery(newData, &descriptor, newProfileData, profileID)

	if newData {
		res, err := database.GetConnection().Collection("user").InsertOne(context.Background(), newProfileData)

		if err != nil {
			return helper.BadRequest("Failed to insert data into the database: " + err.Error())
		}

		profileID = res.InsertedID.(string)

	} else {

		DataUpdateById(newProfileData, profileID)

	}

	res, err := GenerateFromAI(client, aiQuery, "userProfile", &descriptor)
	if err != nil {
		return helper.InternalServerError(err.Error())
	}

	if res == nil {
		return helper.InternalServerError("Open AI not responding")
	} else {
		var embeddingKeywords string
		fmt.Println(res)
		if res["embedding_keywords"] != nil {
			embeddingKeywords = res["embedding_keywords"].(string)
			embeddedData, err := GenerateEmbeddingFromAI(client, embeddingKeywords)
			if err != nil {
				return helper.InternalServerError(err.Error())
			}

			fmt.Println(embeddedData)
			res["embedded"] = embeddedData
		}

	}

	DataUpdateById(res, profileID)

	return helper.SuccessResponse(c, res)
}

func buildAIQuery(newData bool, descriptor *OpenAIDescriptors, newProfileData map[string]interface{}, profileID string) string {
	if newData {

		return fmt.Sprintf(
			"Turn this user profile into a natural language description based on the following user profile model:\n%s\n\n"+
				"ONLY INTERPRET THE DATA IN THE USER PROFILE DATA PROVIDED. DO NOT MAKE UP ANYTHING ELSE.\n\n"+
				"The new user profile data is:\n%s\n\n"+
				"Your task is to:\n"+
				"1. Fill out the following descriptive fields based strictly on the user's data: `my_intro`, `bio`, `professional_journey`, `expertise`, `work_life_philosophy`, and `hobbies`.\n"+
				"2. Extract and return a comma-separated list of important keywords from the entire profile and include it in a field called `embedding_keywords`. These keywords will be used for semantic search and must represent skills, industries, values, and professional traits found in the profile.\n\n"+
				"Return everything as a well-formatted JSON object that includes all the fields above.",
			descriptor.OpenAIDescriptorsConfig().Properties["userProfile"],
			newProfileData,
		)

	}

	existingProfile, err := fetchUserProfile(profileID)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		"Turn this user profile into a natural language description based on the following user profile model:\n%s\n"+
			"ONLY INTERPRET THE DATA IN THE USER PROFILE DATA PROVIDED. DO NOT MAKE UP ANYTHING ELSE.\n"+
			"If geo Location Changed in New Data change geo location of old data also:\n%s\n"+
			"The current user profile data is:\n%s\n"+
			"The new information to add to the user profile is:\n%s\n",
		"Your task is to:\n"+
			"1. Fill out the following descriptive fields based strictly on the user's data: `my_intro`, `bio`, `professional_journey`, `expertise`, `work_life_philosophy`, and `hobbies`.\n"+
			"2. Extract and return a comma-separated list of important keywords from the entire profile and include it in a field called `embedding_keywords`. These keywords will be used for semantic search and must represent skills, industries, values, and professional traits found in the profile.\n\n"+
			"Return everything as a well-formatted JSON object that includes all the fields above.",
		descriptor.OpenAIDescriptorsConfig().Properties["userProfile"], existingProfile, newProfileData,
	)
}

func fetchUserProfile(profileID string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := database.GetConnection().Collection("user").FindOne(context.Background(), bson.M{"_id": profileID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DataUpdateById(data map[string]interface{}, updateId string) {
	filter := bson.M{"_id": updateId}
	update := bson.M{"$set": data}
	database.GetConnection().Collection("user").UpdateOne(context.Background(), filter, update)
}

func MatchUserProfileById(c *fiber.Ctx) error {
	var geoLocation primitive.A
	seen := make(map[interface{}]bool)
	// UserId := c.Params("userId")
	userToken := helper.GetUserTokenValue(c)
	// fmt.Println(userToken)
	userData, err := fetchUserProfile(userToken.UserId)
	if err != nil {

		return helper.BadRequest(err.Error())
	}

	distance := helper.ToInt(c.Params("dis"))
	if distance == 0 {
		distance = 5000
	}

	// var data map[string]interface{}
	// err = c.BodyParser(&data)
	// if err != nil {
	// 	return helper.BadRequest(err.Error())
	// }

	// distance := helper.ToInt(data["distance"])
	if userData["geo_location"] != nil {
		geoLocation = userData["geo_location"].(primitive.A)
	} else {
		return helper.Unexpected("Update Location")
	}

	getUserIds, usersFound := helper.GetSuggestion(userToken.UserId)

	fmt.Println(getUserIds, "userIds")

	var pipeline bson.A
	if usersFound {
		pipeline = bson.A{
			bson.D{
				{"$geoNear",
					bson.D{
						{"near",
							bson.D{
								{"type", "Point"},
								{"coordinates",
									geoLocation,
								},
							},
						},
						{"distanceField", "string"},
						{"maxDistance", distance},
						{"spherical", true},
					},
				},
			},
			bson.D{{"$set", bson.D{{"from_user", userToken.UserId}}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "user_matched"},
						{"localField", "_id"},
						{"foreignField", "user_ids"},
						{"as", "result"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "user"},
						{"localField", "from_user"},
						{"foreignField", "_id"},
						{"as", "from_user_result"},
					},
				},
			},
			bson.D{
				{"$unwind",
					bson.D{
						{"path", "$from_user_result"},
						{"preserveNullAndEmptyArrays", true},
					},
				},
			},
			bson.D{
				{"$addFields",
					bson.D{
						{"userIDS",
							bson.D{
								{"$reduce",
									bson.D{
										{"input", "$result"},
										{"initialValue", bson.A{}},
										{"in",
											bson.D{
												{"$concatArrays",
													bson.A{
														"$$value",
														"$$this.user_ids",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "notifications"},
						{"let",
							bson.D{
								{"from_user", "$from_user"},
								{"to_user", "$_id"},
							},
						},
						{"pipeline",
							bson.A{
								bson.D{
									{"$match",
										bson.D{
											{"$expr",
												bson.D{
													{"$and",
														bson.A{
															bson.D{
																{"$eq",
																	bson.A{
																		"$$from_user",
																		"$from",
																	},
																},
															},
															bson.D{
																{"$eq",
																	bson.A{
																		"$$to_user",
																		"$to",
																	},
																},
															},
															bson.D{
																{"$eq",
																	bson.A{
																		"pending",
																		"$req_status",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{"as", "result"},
					},
				},
			},
			bson.D{
				{"$set",
					bson.D{
						{"request_sent",
							bson.D{
								{"$gt",
									bson.A{
										bson.D{
											{"$size",
												bson.D{
													{"$ifNull",
														bson.A{
															"$result",
															bson.A{},
														},
													},
												},
											},
										},
										0,
									},
								},
							},
						},
						{"from_user_embedded", "$from_user_result.embedded"},
					},
				},
			},
			bson.D{
				{"$match", bson.D{
					{"$and", bson.A{
						bson.D{
							{"userIDS", bson.D{
								{"$nin", bson.A{userToken.UserId}},
							}},
						},
						bson.D{
							{"userIDS", bson.D{
								{"$nin", getUserIds},
							}},
						},
					}},
				}},
			},

			bson.D{
				{"$group",
					bson.D{
						{"_id", primitive.Null{}},
						{"id", bson.D{{"$push", "$_id"}}},
						{"embedded", bson.D{{"$first", "$from_user_embedded"}}},
						{"userData",
							bson.D{
								{"$push",
									bson.D{
										{"user_id", "$_id"},
										{"string", "$string"},
									},
								},
							},
						},
					},
				},
			},
		}
	} else {
		pipeline = bson.A{
			bson.D{
				{"$geoNear",
					bson.D{
						{"near",
							bson.D{
								{"type", "Point"},
								{"coordinates",
									geoLocation,
								},
							},
						},
						{"distanceField", "string"},
						{"maxDistance", distance},
						{"spherical", true},
					},
				},
			},
			bson.D{{"$set", bson.D{{"from_user", userToken.UserId}}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "user_matched"},
						{"localField", "_id"},
						{"foreignField", "user_ids"},
						{"as", "result"},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "user"},
						{"localField", "from_user"},
						{"foreignField", "_id"},
						{"as", "from_user_result"},
					},
				},
			},
			bson.D{
				{"$unwind",
					bson.D{
						{"path", "$from_user_result"},
						{"preserveNullAndEmptyArrays", true},
					},
				},
			},
			bson.D{
				{"$addFields",
					bson.D{
						{"userIDS",
							bson.D{
								{"$reduce",
									bson.D{
										{"input", "$result"},
										{"initialValue", bson.A{}},
										{"in",
											bson.D{
												{"$concatArrays",
													bson.A{
														"$$value",
														"$$this.user_ids",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "notifications"},
						{"let",
							bson.D{
								{"from_user", "$from_user"},
								{"to_user", "$_id"},
							},
						},
						{"pipeline",
							bson.A{
								bson.D{
									{"$match",
										bson.D{
											{"$expr",
												bson.D{
													{"$and",
														bson.A{
															bson.D{
																{"$eq",
																	bson.A{
																		"$$from_user",
																		"$from",
																	},
																},
															},
															bson.D{
																{"$eq",
																	bson.A{
																		"$$to_user",
																		"$to",
																	},
																},
															},
															bson.D{
																{"$eq",
																	bson.A{
																		"pending",
																		"$req_status",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{"as", "result"},
					},
				},
			},
			bson.D{
				{"$set",
					bson.D{
						{"request_sent",
							bson.D{
								{"$gt",
									bson.A{
										bson.D{
											{"$size",
												bson.D{
													{"$ifNull",
														bson.A{
															"$result",
															bson.A{},
														},
													},
												},
											},
										},
										0,
									},
								},
							},
						},
						{"from_user_embedded", "$from_user_result.embedded"},
					},
				},
			},
			bson.D{
				{"$match",
					bson.D{
						{"userIDS",
							bson.D{
								{"$nin",
									bson.A{
										userToken.UserId,
									},
								},
							},
						},
					},
				},
			},

			bson.D{
				{"$group",
					bson.D{
						{"_id", primitive.Null{}},
						{"id", bson.D{{"$push", "$_id"}}},
						{"embedded", bson.D{{"$first", "$from_user_embedded"}}},
						{"userData",
							bson.D{
								{"$push",
									bson.D{
										{"user_id", "$_id"},
										{"string", "$string"},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	// fmt.Println(userToken.UserId)

	var vectorResult []bson.M
	results, _ := helper.GetAggregateQueryResult("user", pipeline)

	if len(results) == 0 {
		var Userresults []bson.M
		Userresults = append(Userresults, userData)
		return helper.SuccessResponse(c, Userresults)
	}

	userIdsResults := results[0]

	if userIdsResults != nil {
		// fmt.Println(userIdsResults)

		userIds := userIdsResults["id"].(primitive.A)
		convertedUserId := helper.ConvertPrimitiveAToStringSlice(userIds)
		helper.SetSuggestion(userToken.UserId, convertedUserId)
		fmt.Println(userIds)
		fromUserEmbeddedData := userIdsResults["embedded"].(primitive.A)
		userDataDis := userIdsResults["userData"].(primitive.A)
		vectorPipeline := bson.A{
			bson.D{
				{"$vectorSearch",
					bson.D{
						{"index", "vector"},
						{"path", "embedded"},
						{"queryVector",
							fromUserEmbeddedData,
						},
						{"numCandidates", 100},
						{"limit", 100},
						{"similarity", "cosine"},
					},
				},
			},
			bson.D{{"$match", bson.D{{"_id", bson.D{{"$in", userIds}}}}}},
			bson.D{
				{"$set",
					bson.D{
						{"score", bson.D{{"$meta", "vectorSearchScore"}}},
						{"string",
							bson.D{
								{"$let",
									bson.D{
										{"vars",
											bson.D{
												{"matched",
													bson.D{
														{"$first",
															bson.D{
																{"$filter",
																	bson.D{
																		{"input", userDataDis},
																		{"as", "item"},
																		{"cond",
																			bson.D{
																				{"$eq",
																					bson.A{
																						"$$item.user_id",
																						bson.D{{"$toString", "$_id"}},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										{"in", "$$matched.string"},
									},
								},
							},
						},
					},
				},
			},
			bson.D{{"$set", bson.D{{"from_user", userToken.UserId}}}},
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "notifications"},
						{"let",
							bson.D{
								{"from_user", "$from_user"},
								{"to_user", "$_id"},
							},
						},
						{"pipeline",
							bson.A{
								bson.D{
									{"$match",
										bson.D{
											{"$expr",
												bson.D{
													{"$and",
														bson.A{
															bson.D{
																{"$eq",
																	bson.A{
																		"$$from_user",
																		"$from",
																	},
																},
															},
															bson.D{
																{"$eq",
																	bson.A{
																		"$$to_user",
																		"$to",
																	},
																},
															},
															bson.D{
																{"$eq",
																	bson.A{
																		"pending",
																		"$req_status",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{"as", "result"},
					},
				},
			},
			bson.D{
				{"$set",
					bson.D{
						{"request_sent",
							bson.D{
								{"$gt",
									bson.A{
										bson.D{
											{"$size",
												bson.D{
													{"$ifNull",
														bson.A{
															"$result",
															bson.A{},
														},
													},
												},
											},
										},
										0,
									},
								},
							},
						},
						{"from_user_embedded", "$from_user_result.embedded"},
					},
				},
			},
		}
		// fmt.Println(vectorPipeline)
		vectorResult, _ = helper.GetAggregateQueryResult("user", vectorPipeline)
	}

	fmt.Println(len(vectorResult))
	// Create a map to track seen user IDs
	// seen := make(map[interface{}]bool)

	// Create a new result slice
	var uniqueResults []bson.M

	// Iterate through vectorResult to remove existing duplicates
	for _, user := range vectorResult {
		if userID, ok := user["_id"]; ok {
			if !seen[userID] {
				seen[userID] = true
				uniqueResults = append(uniqueResults, user)
			}
		}
	}

	// Now check userData
	if userID, ok := userData["_id"]; ok {
		if !seen[userID] {
			seen[userID] = true
			uniqueResults = append(uniqueResults, userData)
		}
	}

	return helper.SuccessResponse(c, uniqueResults)
	// fmt.Println(results, "   ", len(results))
	var strcut = OpenAIDescriptors{}

	aiQuery := fmt.Sprintf(
		"You are an assistant that compares users. Return users from 'otherUsers' who have at least one matching attribute with the 'primaryUser'. "+
			"For each matched user, return a match_type and a match_reason. "+
			"The match_reason must be exactly three lines long, providing a clear explanation of why the match was made."+
			"STRICTLY FOLLOW THIS RULE.\n\n"+
			"Primary user profile:\n%v\n\n"+
			"List of other users to check for matches:\n%v\n\n"+
			"Output format (no extra text): ",
		userData, vectorResult,
	)

	config := openai.DefaultConfig(openAiKey)
	config.OrgID = openAiOrgId

	// var res map[string]interface{}
	client := openai.NewClientWithConfig(config)
	res, err := ProfileMatchFromOpenAI(client, aiQuery, "match_reason", &strcut)

	if err != nil {
		// return nil
		return helper.InternalServerError(err.Error())
	}

	// return helper.SuccessResponse(c, res)
	if res == nil {
		var userIds []interface{}
		// userIds = append(userIds, userToken.UserId)
		newRes := map[string]interface{}{
			"user_id":      userToken.UserId,
			"match_reason": "Proficiency in Angular and Node.js with experience in industry-relevant technologies like Go, sharing professional focus and skills.",
			"match_type":   "Fit",
		}
		userIds = append(userIds, newRes)

		res["matched_users"] = userIds
	}

	matchedUserData := res["matched_users"].([]interface{})

	// var Userresults []bson.M
	// Track unique _id values

	// Iterate over user IDs and filter results
	for _, users := range matchedUserData {
		userMap := users.(map[string]interface{})
		for _, user := range vectorResult {
			userId := userMap["user_id"].(string)
			if user["_id"] == userId {
				if !seen[userId] { // Check if already added
					user["match_reason"] = userMap["match_reason"].(string)
					user["match_type"] = userMap["match_type"].(string)
					vectorResult = append(vectorResult, user)
					seen[userId] = true
				}
			}
		}
	}

	// Add userData if it's not already in Userresults
	if userID, exists := userData["_id"]; exists {
		if !seen[userID] { // Check if userData is unique
			vectorResult = append(vectorResult, userData)
			seen[userID] = true
		}
	}

	return helper.SuccessResponse(c, vectorResult)
}

func MatchAllUserProfile(c *fiber.Ctx) error {
	// var geoLocation primitive.A
	// UserId := c.Params("userId")
	userToken := helper.GetUserTokenValue(c)
	fmt.Println(userToken)
	userData, err := fetchUserProfile(userToken.UserId)
	if err != nil {
		return helper.BadRequest(err.Error())
	}

	// var data map[string]interface{}
	// err = c.BodyParser(&data)
	// if err != nil {
	// 	return helper.BadRequest(err.Error())
	// }

	// pipeline := bson.A{
	// 	bson.D{
	// 		{"$geoNear",
	// 			bson.D{
	// 				{"near",
	// 					bson.D{
	// 						{"type", "Point"},
	// 						{"coordinates",
	// 							geoLocation,
	// 						},
	// 					},
	// 				},
	// 				{"distanceField", "string"},
	// 				{"maxDistance", 50000},
	// 				{"spherical", true},
	// 			},
	// 		},
	// 	},
	// 	// bson.D{{"$match", bson.D{{"_id", bson.D{{"$ne", inputData.UserId}}}}}},
	// }

	// pipeline := bson.A{
	// 	bson.D{
	// 		{"$geoNear",
	// 			bson.D{
	// 				{"near",
	// 					bson.D{
	// 						{"type", "Point"},
	// 						{"coordinates",
	// 							geoLocation,
	// 						},
	// 					},
	// 				},
	// 				{"distanceField", "string"},
	// 				{"maxDistance", distance},
	// 				{"spherical", true},
	// 			},
	// 		},
	// 	},
	// 	bson.D{
	// 		{"$lookup",
	// 			bson.D{
	// 				{"from", "user_matched"},
	// 				{"localField", "_id"},
	// 				{"foreignField", "user_ids"},
	// 				{"as", "result"},
	// 			},
	// 		},
	// 	},
	// 	bson.D{
	// 		{"$addFields",
	// 			bson.D{
	// 				{"userIDS",
	// 					bson.D{
	// 						{"$reduce",
	// 							bson.D{
	// 								{"input", "$result"},
	// 								{"initialValue", bson.A{}},
	// 								{"in",
	// 									bson.D{
	// 										{"$concatArrays",
	// 											bson.A{
	// 												"$$value",
	// 												"$$this.user_ids",
	// 											},
	// 										},
	// 									},
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	bson.D{
	// 		{"$match",
	// 			bson.D{
	// 				{"result.user_ids",
	// 					bson.D{
	// 						{"$nin",
	// 							bson.A{
	// 								userToken.UserId,
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	pipeline := bson.A{

		bson.D{{"$set", bson.D{{"from_user", userToken.UserId}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "user_matched"},
					{"localField", "_id"},
					{"foreignField", "user_ids"},
					{"as", "result"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"userIDS",
						bson.D{
							{"$reduce",
								bson.D{
									{"input", "$result"},
									{"initialValue", bson.A{}},
									{"in",
										bson.D{
											{"$concatArrays",
												bson.A{
													"$$value",
													"$$this.user_ids",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "notifications"},
					{"let",
						bson.D{
							{"from_user", "$from_user"},
							{"to_user", "$_id"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$and",
													bson.A{
														bson.D{
															{"$eq",
																bson.A{
																	"$$from_user",
																	"$from",
																},
															},
														},
														bson.D{
															{"$eq",
																bson.A{
																	"$$to_user",
																	"$to",
																},
															},
														},
														bson.D{
															{"$eq",
																bson.A{
																	"pending",
																	"$req_status",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{"as", "result"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"request_sent",
						bson.D{
							{"$gt",
								bson.A{
									bson.D{
										{"$size",
											bson.D{
												{"$ifNull",
													bson.A{
														"$result",
														bson.A{},
													},
												},
											},
										},
									},
									0,
								},
							},
						},
					},
				},
			},
		},
		// bson.D{
		// 	{"$match",
		// 		bson.D{
		// 			{"result.user_ids",
		// 				bson.D{
		// 					{"$nin",
		// 						bson.A{
		// 							userToken.UserId,
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		bson.D{
			{"$match",
				bson.D{
					{"userIDS",
						bson.D{
							{"$nin",
								bson.A{
									userToken.UserId,
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println(pipeline)

	results, err := helper.GetAggregateQueryResult("user", pipeline)
	// if err != nil {
	// 	return helper.BadRequest(err.Error())
	// }

	if len(results) == 0 {
		var Userresults []bson.M
		Userresults = append(Userresults, userData)
		return helper.SuccessResponse(c, Userresults)
	}
	return helper.SuccessResponse(c, results)
	fmt.Println(results, "   ", len(results))
	var strcut = OpenAIDescriptors{}
	aiQuery := fmt.Sprintf(
		"You are a strict filtering engine. Your task is to find and return ONLY the users from the following list whose key_skills ,expertise have at least ONE matches with the primary user’s key_skills ,expertise. "+
			"DO NOT return users with fewer than ONE common key_skills ,expertise. "+
			"If no users match this condition, return an empty array. "+
			"STRICTLY FOLLOW THIS RULE.\n\n"+
			"Primary user profile:\n%v\n\n"+
			"List of other users to check for matches:\n%v\n\n"+
			"Output format (no extra text): ",
		userData, results,
	)

	config := openai.DefaultConfig(openAiKey)
	config.OrgID = openAiOrgId

	client := openai.NewClientWithConfig(config)

	res, err := ProfileMatchFromOpenAI(client, aiQuery, "match_profile", &strcut)

	if err != nil {
		// return nil
		return helper.InternalServerError(err.Error())
	}

	// fmt.Println(res)

	// UserIds := res["user_ids"].([]interface{})

	var Userresults []bson.M
	// for _, id := range UserIds {
	// 	for _, user := range results {
	// 		if user["_id"] == id {
	// 			Userresults = append(Userresults, user)
	// 		}
	// 	}
	// }

	if res == nil {
		var userIds []interface{}
		// userIds = append(userIds, userToken.UserId)
		newRes := map[string]interface{}{
			"user_id":      userToken.UserId,
			"match_reason": "Proficiency in Angular and Node.js with experience in industry-relevant technologies like Go, sharing professional focus and skills.",
			"match_type":   "Fit",
		}
		userIds = append(userIds, newRes)

		res["matched_users"] = userIds
	}

	matchedUserData := res["matched_users"].([]interface{})

	// var Userresults []bson.M
	seen := make(map[interface{}]bool) // Track unique _id values

	// Iterate over user IDs and filter results
	for _, users := range matchedUserData {
		userMap := users.(map[string]interface{})
		for _, user := range results {
			userId := userMap["user_id"].(string)
			if user["_id"] == userId {
				if !seen[userId] { // Check if already added
					user["match_reason"] = userMap["match_reason"].(string)
					user["match_type"] = userMap["match_type"].(string)
					Userresults = append(Userresults, user)
					seen[userId] = true
				}
			}
		}
	}

	// Add userData if it's not already in Userresults
	if userID, exists := userData["_id"]; exists {
		if !seen[userID] { // Check if userData is unique
			Userresults = append(Userresults, userData)
			seen[userID] = true
		}
	}

	Userresults = append(Userresults, userData)

	return helper.SuccessResponse(c, Userresults)
}

func UpdateProfileById(c *fiber.Ctx) error {
	profileID := c.Params("profileId")
	var newProfileData map[string]interface{}

	if err := c.BodyParser(&newProfileData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseModel{
			Status:  "failure",
			Message: "Error parsing request body",
			Details: err.Error(),
		})
	}

	newData := profileID == ""

	existingProfile1, _ := fetchUserProfile(profileID)
	if existingProfile1 == nil {
		newData = true
	}

	client := openai.NewClient(openAiKey)

	var descriptor OpenAIDescriptors
	aiQuery := buildAIQuery(newData, &descriptor, newProfileData, profileID)

	if newData {
		res, err := database.GetConnection().Collection("user").InsertOne(context.Background(), newProfileData)

		if err != nil {
			return helper.BadRequest("Failed to insert data into the database: " + err.Error())
		}

		profileID = res.InsertedID.(string)

	} else {

		DataUpdateById(newProfileData, profileID)
	}

	res, err := GenerateFromAI(client, aiQuery, "userProfile", &descriptor)
	if err != nil {
		return helper.InternalServerError(err.Error())
	}

	if res == nil {
		return helper.InternalServerError("Open AI not responding")
	}

	DataUpdateById(res, profileID)

	return helper.SuccessResponse(c, res)
}

func GetUserOnboardingController(c *fiber.Ctx) error {

	userId := c.Params("userId")
	var aiQuery string
	var strcut = OpenAIDescriptors{}
	questionLimit := 10
	if userId == "" {
		aiQuery = fmt.Sprintf(
			"Return a list of questions to enhance the existing user's profile. "+
				"The aim is to enhance information about the user to optimize their matching with other professional profiles. "+
				"Ensure this is based on dating-type matching but for professionals. \n"+
				"The structure of these questions should aim to complete the user data model: \n"+
				"%s"+
				"The question limit is a maximum of %d questions \n",
			// Assuming OpenAIDescriptorsConfig().Properties["userProfile"] is serialized to a string
			strcut.OpenAIDescriptorsConfig().Properties["onboardingQuestions"], questionLimit,
		)
	} else {
		var userData map[string]interface{}
		err := database.GetConnection().Collection("user").FindOne(context.Background(), bson.M{"_id": userId}).Decode(&userData)
		if err != nil {
			return helper.BadRequest("User Not Found")
		}
		aiQuery = fmt.Sprintf(
			"Return a list of questions to enhance the existing user's profile.  \n"+
				"%s"+
				"Ask Only the remaining new question . \n"+
				"The current user profile data is: \n"+"%s"+"\n"+
				"",
			// Assuming OpenAIDescriptorsConfig().Properties["userProfile"] is serialized to a string
			strcut.OpenAIDescriptorsConfig().Properties["onboardingQuestions"], userData,
		)
	}

	config := openai.DefaultConfig(openAiKey)
	config.OrgID = openAiOrgId

	client := openai.NewClientWithConfig(config)

	// Construct the AI query string using fmt.Sprintf

	res, err := GenerateFromAI(client, aiQuery, "onboardingQuestions", &strcut)
	if err != nil {
		return helper.InternalServerError(err.Error())
	}

	return helper.SuccessResponse(c, res)
}
