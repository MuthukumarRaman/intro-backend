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
)

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
	Description string              `json:"description,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Default     string              `json:"default,omitempty"`
}

func (d *OpenAIDescriptors) OpenAIDescriptorsConfig() Property {
	return Property{
		Type: "object",
		Properties: map[string]Property{
			// "userProfile": {
			// 	Type:        "object",
			// 	Description: "A natural description of who you are as a person",
			// 	Properties: map[string]Property{
			// 		"name": {
			// 			Type:        "string",
			// 			Description: "The name of the user",
			// 		},
			// 		"age": {
			// 			Type:        "number",
			// 			Description: "The age of the user",
			// 		},
			// 		"location": {
			// 			Type:        "string",
			// 			Description: "The location of the user",
			// 		},
			// 		"industry": {
			// 			Type:        "string",
			// 			Description: "The industry of the user",
			// 		},
			// 		"introduction": {
			// 			Type:        "string",
			// 			Description: "A brief introduction about yourself, including your name, age, and where you're based",
			// 		},
			// 		"professionalStory": {
			// 			Type:        "string",
			// 			Description: "Tell us about your professional journey, current role, and what you're passionate about in your work",
			// 		},
			// 		"expertise": {
			// 			Type:        "string",
			// 			Description: "What are your key areas of expertise and skills?",
			// 		},
			// 		"personalInterests": {
			// 			Type:        "string",
			// 			Description: "Share what you love to do outside of work - your hobbies, interests, and what makes you unique",
			// 		},
			// 		"workStyle": {
			// 			Type:        "string",
			// 			Description: "Describe how you prefer to work and collaborate with others",
			// 		},
			// 		"aspirations": {
			// 			Type:        "string",
			// 			Description: "What are your goals and what kind of opportunities or connections are you looking for?",
			// 		},
			// 		"funFact": {
			// 			Type:        "string",
			// 			Description: "Share an interesting story or fact about yourself that helps people remember you",
			// 		},
			// 	},
			// },
			// "onboardingQuestions": {
			// 	Type:        "object",
			// 	Description: "Return an onboarding guide as natural language using a structured JSON model. Each question should have a `question` field for the prompt text and a `type` field for the input type (e.g., input, date, textarea, radio).",
			// 	Properties: map[string]Property{
			// 		"questions": {
			// 			Type:        "array",
			// 			Description: "A list of onboarding questions with their details.",
			// 			Items: &Property{
			// 				Type:        "object",
			// 				Description: "A natural language question to ask the user",
			// 				Properties: map[string]Property{
			// 					"question": {
			// 						Type:        "string",
			// 						Description: "The text of the question to ask the user.",
			// 					},
			// 					"type": {
			// 						Type:        "string",
			// 						Description: "The type of input for the question (e.g., input, date, textarea, radio).",
			// 					},
			// 					"sampleAnswer": {
			// 						Type:        "string",
			// 						Description: "An example answer for the question.",
			// 					},
			// 					"dataType": {
			// 						Type: "string",
			// 						// enum: ["string", "number", "date", "boolean"],
			// 						Description: "The type of the answer (string, number, date, or etc).",
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },

			// "userProfile": {
			// 	Type:        "object",
			// 	Description: "A natural description of who you are as a person and give me as a single object",
			// 	Properties: map[string]Property{
			// 		"name": {
			// 			Type:        "string",
			// 			Description: "The name of the user",
			// 		},
			// 		"date_of_birth": {
			// 			Type:        "string",
			// 			Format:      "date-time",
			// 			Description: "The date of birth of the user in ISO 8601 date-time format (e.g., YYYY-MM-DD)",
			// 		},
			// 		"age": {
			// 			Type:        "number",
			// 			Description: "Calculate the age by using given date of birth",
			// 		},
			// 		"location": {
			// 			Type:        "string",
			// 			Description: "The location of the user",
			// 		},
			// 		"industry": {
			// 			Type:        "string",
			// 			Description: "The industry of the user",
			// 		},
			// 		"introduction": {
			// 			Type:        "string",
			// 			Description: "A brief introduction about yourself, including your name, age, and where you're based",
			// 		},
			// 		"professional_story": {
			// 			Type:        "string",
			// 			Description: "Tell us about your professional journey, current role, and what you're passionate about in your work",
			// 		},
			// 		"expertise": {
			// 			Type:        "string",
			// 			Description: "What are your key areas of expertise and skills?",
			// 		},
			// 		"personal_interests": {
			// 			Type:        "string",
			// 			Description: "Share what you love to do outside of work - your hobbies, interests, and what makes you unique",
			// 		},
			// 		"work_style": {
			// 			Type:        "string",
			// 			Description: "Describe how you prefer to work and collaborate with others",
			// 		},
			// 		"aspirations": {
			// 			Type:        "string",
			// 			Description: "What are your goals and what kind of opportunities or connections are you looking for?",
			// 		},
			// 		"fun_fact": {
			// 			Type:        "string",
			// 			Description: "Share an interesting story or fact about yourself that helps people remember you",
			// 		},
			// 	},
			// },

			"userProfile": {
				Type:        "object",
				Description: "A comprehensive profile where you describe yourself.",
				Required:    []string{"my_intro", "bio", "age", "location", "industry", "work_life_philosophy", "professional_journey", "expertise", "hobbies"},
				Properties: map[string]Property{
					"my_intro": {
						Type:        "string",
						Description: "Write a short introduction about yourself.",
						Default:     "I am an enthusiastic professional eager to learn and grow.",
					},
					"bio": {
						Type:        "string",
						Description: "Describe yourself briefly, including your profession and experience.",
						Default:     "I am a dedicated professional with a strong background in my field.",
					},
					"age": {
						Type:        "number",
						Description: "Enter your age, calculated from your date of birth.",
						Default:     "25",
					},
					"location": {
						Type:        "string",
						Description: "Mention your current location (city and country).",
						Default:     "Not specified",
					},
					"industry": {
						Type:        "string",
						Description: "List the industries you have experience in.",
						Default:     "Technology",
					},
					"work_life_philosophy": {
						Type:        "string",
						Description: "Share your thoughts on work-life balance and your professional approach.",
						Default:     "I believe in maintaining a healthy balance between work and personal life.",
					},
					"professional_journey": {
						Type:        "string",
						Description: "Summarize your career path and what led you to your current role.",
						Default:     "I started my career with a passion for problem-solving and grew into my current role through continuous learning.",
					},
					"expertise": {
						Type:        "string",
						Description: "List your key skills and areas of expertise.",
						Default:     "Problem-solving, Communication, Technical Skills",
					},
					"hobbies": {
						Type:        "string",
						Description: "Mention your hobbies and activities outside of work.",
						Default:     "Reading, Traveling, Fitness",
					},
				},
			},
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
			"userProfile": {
				ModelID:        "gpt-3.5-turbo",
				TaskDefinition: "Return a user profile as natural language using a JSON structure.",
				AIFunctions: []AIFunction{
					{
						Name:        "parseToUserProfileModel",
						Description: "Parse a user profile description to a JSON model",
						Parameters:  descriptors.OpenAIDescriptorsConfig(),
					},
				},
				ForcedFunction: "parseToUserProfileModel",
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

type ResponseModel struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Details string      `json:"details"`
	Data    interface{} `json:"data,omitempty"`
}

func GetUserProfile(c *fiber.Ctx) error {

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

	// Call the controller
	if profileID != "" {
		// var err error
		// result, err = UserProfileUpdateController(profileID, newProfileData)
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(ResponseModel{
		// 		Status:  "failure",
		// 		Message: "Internal server error",
		// 		Details: err.Error(),
		// 	})
		// }
		DataUpdateById(newProfileData, profileID)
	}

	client := openai.NewClient("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")
	var strcut = OpenAIDescriptors{}
	var aiQuery string

	if result == nil {
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

func DataUpdateById(data map[string]interface{}, updateId string) {

	filter := bson.M{
		"_id": updateId,
	}

	update := bson.M{
		"$set": data,
	}

	database.GetConnection().Collection("user").UpdateOne(context.Background(), filter, update)
}
