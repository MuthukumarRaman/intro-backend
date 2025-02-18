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

	fmt.Println(resp)

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

	client := openai.NewClient("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")
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

// func DataUpdateById(data map[string]interface{}, updateId string) {

// 	filter := bson.M{
// 		"_id": updateId,
// 	}

// 	update := bson.M{
// 		"$set": data,
// 	}

// 	database.GetConnection().Collection("user").UpdateOne(context.Background(), filter, update)
// }

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

	client := openai.NewClient("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")

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

	if res["userProfile"] == nil {
		return helper.InternalServerError("Open AI not responding")
	}

	DataUpdateById(res["userProfile"].(map[string]interface{}), profileID)

	return helper.SuccessResponse(c, res)
}

func buildAIQuery(newData bool, descriptor *OpenAIDescriptors, newProfileData map[string]interface{}, profileID string) string {
	if newData {
		return fmt.Sprintf(
			"Turn this user profile into a natural language description based on the following user profile model:\n%s\n"+
				"ONLY INTERPRET THE DATA IN THE USER PROFILE DATA PROVIDED. DO NOT MAKE UP ANYTHING ELSE.\n"+
				"The new information to add to the user profile is:\n%s\n",
			descriptor.OpenAIDescriptorsConfig().Properties["userProfile"], newProfileData,
		)
	}

	existingProfile, err := fetchUserProfile(profileID)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		"Turn this user profile into a natural language description based on the following user profile model:\n%s\n"+
			"ONLY INTERPRET THE DATA IN THE USER PROFILE DATA PROVIDED. DO NOT MAKE UP ANYTHING ELSE.\n"+
			"The current user profile data is:\n%s\n"+
			"The new information to add to the user profile is:\n%s\n",
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
