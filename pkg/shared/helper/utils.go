package helper

import (
	// "errors"
	// "fmt"

	"context"
	"fmt"

	"reflect"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"introme-api/pkg/shared/database"
	// "introme-api/pkg/shared/helper"
)

type StatusError struct {
	Code       int
	Message    string
	Validation interface{}
}

// Error implements error.

func UpdateDateObject(input map[string]interface{}) error {
	for k, v := range input {
		if v == nil {
			continue
		}
		ty := reflect.TypeOf(v).Kind().String()
		if ty == "string" {
			val := reflect.ValueOf(v).String()
			t, err := time.Parse(time.RFC3339, val)
			if err == nil {
				input[k] = t.UTC()
			}
		} else if ty == "map" {
			return UpdateDateObject(v.(map[string]interface{}))
		} else if ty == "slice" {
			for _, e := range v.([]interface{}) {
				if reflect.TypeOf(e).Kind().String() == "map" {
					return UpdateDateObject(e.(map[string]interface{}))
				}
			}
		}
	}
	return nil
}

func Toint64(s string) int64 {
	if s == "" {
		return int64(0)
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func Page(s string) int64 {
	return Toint64(s)
}

func Limit(s string) int64 {
	if s == "" {
		s = GetenvStr("DEFAULT_FETCH_ROWS")
	}
	return Toint64(s)
}

func DocIdFilter(id string) bson.M {
	if id == "" {
		return bson.M{}
	}
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return bson.M{"_id": id}
	} else {
		return bson.M{"_id": docId}
	}
}

func ObjectIdToString(id interface{}) string {
	return id.(primitive.ObjectID).Hex()
}

func PostDataModelConfig(c *fiber.Ctx) error {
	// Get parameters and headers
	collectionName := c.Params("model_name")
	orgId := c.Get("OrgId")

	if orgId == "" {
		return BadRequest("Organization Id missing")
	}

	// Initialize variables
	var insertData interface{}
	users := GetUserTokenValue(c)

	switch collectionName {
	case "model_config":
		var modelData model_config
		if err := c.BodyParser(&modelData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}
		modelData.Status = "A"
		modelData.CreatedBy = users.UserId
		modelData.CreatedOn = time.Now()
		insertData = modelData
	case "data_model":
		var configData Config
		if err := c.BodyParser(&configData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}
		configData.Status = "A"
		insertData = configData
	case "user":
		return SendSimpleEmailHandler(c)
	case "screen":
		var screen Screen
		if err := c.BodyParser(&screen); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}
		screen.CreatedBy = users.UserId
		screen.CreatedOn = time.Now()
		screen.Status = "A"
		insertData = screen

	case "group":
		var grouping PaginationRequest

		if err := c.BodyParser(&grouping); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing request body")
		}
		grouping.CreatedBy = users.UserId
		grouping.CreatedOn = time.Now()
		grouping.Status = "A"
		insertData = grouping
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Unsupported collection name",
		})
	}

	// Insert the data into the database
	res, err := database.GetConnection().Collection(collectionName).InsertOne(ctx, insertData)
	if err != nil {

	}

	responseJSON := fiber.Map{
		"status":  "success",
		"message": "Data Added Successfully",
		"data": fiber.Map{
			"InsertedID": res.InsertedID,
		},
	}
	return c.Status(fiber.StatusOK).JSON(responseJSON)
}

func UpdateData(c *fiber.Ctx, collectionName string, filter interface{}, data map[string]interface{}, orgId string) error {

	// // Validate the input data based on the data model
	// inputData, validationErrors := UpdateValidateInDatamodel(collectionName, data, orgId)
	// if validationErrors != nil {
	// 	// Handle validation errors with status code 400 (Bad Request)
	// 	jsonstring, _ := json.Marshal(validationErrors)
	// 	return BadRequest(string(jsonstring))
	// }

	// updatedData := make(map[string]interface{})
	// Data := updateFieldsWithParentKey(inputData, "", updatedData)

	update := bson.M{
		"$set": data,
	}
	fmt.Println(filter, collectionName)

	// Update data in the collection
	_, err := database.GetConnection().Collection(collectionName).UpdateOne(context.Background(), filter, update)
	if err != nil {
		// Handle database update error with status code 500 (Internal Server Error)
		return BadRequest(err.Error())
	}

	return SuccessResponse(c, "Updated Successfully")
}

func GetNextSeqNumber(key string) int32 {
	//update to database
	filter := bson.M{"_id": key}
	updateData := bson.M{
		"$inc": bson.M{"value": 1},
	}
	result, _ := ExecuteFindAndModifyQuery("sequence", filter, updateData)
	return result["value"].(int32)
}

func ToString(input interface{}) string {
	return fmt.Sprintf("%v", input)
}
