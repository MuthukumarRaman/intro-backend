package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	// "strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"introme-api/pkg/shared/database"
)

var links = make(map[string]*Link)

// Search EntitiesHandler - Get Entities
// todo lookup
// func SearchEntityWithChildCountHandler(c *fiber.Ctx) error {
// 	orgId := c.Get("OrgId")
// 	if orgId == "" {
// 		return BadRequest("Organization Id missing")
// 	}
// 	var parentCollection = c.Params("parent_collection")
// 	var keyColumn = c.Params("key_column")
// 	var childCollection = c.Params("child_collection")
// 	var lookupColumn = c.Params("lookup_column")

// 	var conditions []Filter
// 	err := c.BodyParser(&conditions)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}

// 	response, err := GetSearchQueryWithChildCount(orgId, parentCollection, keyColumn, childCollection, lookupColumn, conditions)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	return SuccessResponse(c, response)
// }

// Search EntitiesHandler - Get Entities
// func DataLookupDocsHandler(c *fiber.Ctx) error {
// 	orgId := c.Get("OrgId")
// 	if orgId == "" {
// 		return BadRequest("Organization Id missing")
// 	}
// 	var lookupQuery LookupQuery
// 	err := c.BodyParser(&lookupQuery)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	response, err := ExecuteLookupQuery(orgId, lookupQuery)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	return SuccessResponse(c, response)
// }

func GetDeviceDataByOrganizationID(c *fiber.Ctx) error {
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}

	org_id := c.Params("organisation")

	/*
	   match condtion  key and value
	   lookup
	   from collection field

	*/

	filter := bson.A{
		bson.D{{"$match", bson.D{{"org_id", org_id}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "device"},
					{"localField", "_id"},
					{"foreignField", "facility_id"},
					{"as", "device"},
				},
			},
		},
		bson.D{{"$match", bson.D{{"device.status", "Active"}}}},
	}

	response, err := GetAggregateQueryResult("facility", filter)
	if err != nil {
		return BadRequest(err.Error())
	}
	return SuccessResponse(c, response)

}

// func GrpReportHandler(c *fiber.Ctx) error {
// 	orgId := c.Get("OrgId")
// 	if orgId == "" {
// 		return BadRequest("Organization Id missing")
// 	}
// 	var request GroupSumRequest
// 	err := c.BodyParser(&request)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	response, err := ExecuteGroupReportQuery(orgId, request)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	return SuccessResponse(c, response)
// }

func sharedDBEntityHandler(c *fiber.Ctx) error {
	var collectionName = c.Params("collectionName")
	if collectionName == "db_config" {
		return BadRequest("Access Denied")
	}
	cur, err := database.SharedDB.Collection(collectionName).Find(ctx, bson.D{})
	if err != nil {
		return BadRequest(err.Error())
	}
	var response []bson.M
	if err = cur.All(ctx, &response); err != nil {
		return BadRequest(err.Error())
	}
	return SuccessResponse(c, response)
}

// func ReportHandler(c *fiber.Ctx) error {
// 	orgId := c.Get("OrgId")
// 	if orgId == "" {
// 		return BadRequest("Organization Id missing")
// 	}
// 	var request GroupSumRequest
// 	err := c.BodyParser(&request)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	response, err := ExecuteGroupQuery(orgId, request)
// 	if err != nil {
// 		return BadRequest(err.Error())
// 	}
// 	return SuccessResponse(c, response)
// }

func GenerateAppaccesscode() string {
	// Generate a UUID (Universally Unique Identifier).
	uuidObj := uuid.New()
	return uuidObj.String()
}

func Triggerapi(c *fiber.Ctx) error {
	// Retrieve the parameters from the query string
	// param1 := c.Query("email")
	param2 := c.Query("decoding")

	// Use the parameters as needed in your backend logic
	// For example, you can return them as part of the response
	// response := map[string]string{
	// 	"email":    param1,
	// 	"decoding": param2,
	// }
	links := Accesskeychecking(param2)
	return c.Redirect(links, fiber.StatusMovedPermanently)
	// CheckTempCollection(c,"amsort", param1, param2)
	// ChecktheuserDetails(c,param1, param2)
	// fmt.Println("CHECK THE OUTPUT",response)
	// return nil
}

func Accesskeychecking(accesskey string) string {

	link := fmt.Sprintf("http://localhost:4200/activate?accesskey=%s", accesskey)

	return link

}

func UpdateUserPasswordAndDeleteTempData(c *fiber.Ctx) error {

	var inputData map[string]interface{}
	err := c.BodyParser(&inputData)
	if err != nil {
		return BadRequest("Error parsing request body: " + err.Error())
	}
	access_key := c.Params("access_key")

	query := bson.M{"access_key": access_key}
	//var response []primitive.M
	response, err := FindDocs("amsort", "temporary_user", query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve user", "error": err.Error()})
	}

	ID, idExists := response["_id"].(string)
	if !idExists {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Invalid response format"})
	}

	passwordHash, _ := GeneratePasswordHash(inputData["password"].(string))
	delete(inputData, "password")
	delete(inputData, "confirm_password")
	update := bson.M{"$set": bson.M{"pwd": passwordHash}}
	filter := bson.M{"_id": ID}

	result, err := ExecuteFindAndModifyQuery("user", filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update document", "error": err.Error()})
	}
	_, err = ExecuteDeleteManyByIds("amsort", "temporary_user", filter)
	if err != nil {
		return BadRequest("Failed to Delete data into the database: " + err.Error())

	}

	return SuccessResponse(c, result)
}

func AccessLinkHandler(c *fiber.Ctx) error {
	token := c.Params("token")
	link, ok := links[token]

	if !ok {
		return c.SendStatus(http.StatusNotFound)
	}

	if time.Now().After(link.Expiration) {
		return c.Status(http.StatusGone).JSON(fiber.Map{
			"message": "Your session has expired. Please generate a new link.",
		})
	}

	return c.Redirect(link.URL)
}

func User_junked_files(requestmail string, apptoken string) error {
	requestData := make(map[string]interface{})
	requestData["_id"] = requestmail
	requestData["access_key"] = apptoken
	requestData["expire_on"] = time.Now()

	_, err := database.GetConnection().Collection("temporary_user").InsertOne(ctx, requestData)
	if err != nil {
		return BadRequest("Failed to insert data into the database: " + err.Error())
	}

	return nil
}

func InsertValidateInDatamodel(collectionName, inputJsonString, orgId string) (map[string]interface{}, map[string]string) {
	var validationErrors = make(map[string]string)

	newStructValue, errorMessage := CreateInstanceForCollection(collectionName)
	if len(errorMessage) > 0 {
		return nil, errorMessage
	}

	err := json.Unmarshal([]byte(inputJsonString), newStructValue)
	if err != nil {
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			expectedType := unmarshalErr.Type.String()
			dataType := strings.TrimPrefix(expectedType, "*")
			fieldName := unmarshalErr.Field
			return nil, map[string]string{
				"field":             fieldName,
				"Expected DataType": dataType,
			}
		}

		return nil, nil
	}

	// loop through pointer to get the actual struct
	rv := reflect.ValueOf(newStructValue)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal([]byte(inputJsonString), &inputMap); err != nil {
		return nil, map[string]string{"error": "Invalid JSON data: " + err.Error()}
	}
	//Check the field any extra field is here
	if err := vertifyInputStruct(rv, inputMap, validationErrors); err != nil {
		return nil, validationErrors
	}
	// fmt.Println(rv.Interface())
	validationErrors = ValidateStruct(rv.Interface())
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	// fmt.Println(validationErrors)
	var cleanedData map[string]interface{}
	inputByte, _ := json.Marshal(rv.Interface())
	json.Unmarshal(inputByte, &cleanedData)

	return cleanedData, nil
}

func UpdateValidateInDatamodel(collectionName string, inputJsonString, orgId string) (map[string]interface{}, map[string]string) {
	// newStructFields := DynamicallyBindStructOnDataModel(collectionName, orgId)
	newStructFields, errorMessage := CreateInstanceForCollection(collectionName)
	if len(errorMessage) > 0 {
		return nil, errorMessage
	}

	err := json.Unmarshal([]byte(inputJsonString), &newStructFields)
	if err != nil {
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			expectedType := unmarshalErr.Type.String()
			dataType := strings.TrimPrefix(expectedType, "*")
			fieldName := unmarshalErr.Field
			return nil, map[string]string{
				"field":             fieldName,
				"Expected DataType": dataType,
			}
		}
		return nil, map[string]string{"error": "Failed to unmarshal input JSON"}
	}

	// loop through pointer to get the actual struct
	rv := reflect.ValueOf(newStructFields)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal([]byte(inputJsonString), &inputMap); err != nil {
		return nil, map[string]string{"error": "Invalid JSON data: " + err.Error()}
	}

	matchedFields := FilterStructFieldsByJSON(rv, inputMap)
	newStructType := reflect.StructOf(matchedFields)

	// Create a new struct instance with the matched fields
	newStruct := reflect.New(newStructType).Interface()

	err = json.Unmarshal([]byte(inputJsonString), &newStruct)
	if err != nil {
		// return nil, map[string]string{"error": "Failed to unmarshal input JSON"}
	}

	validationErrors := ValidateStruct(newStruct)
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	var cleanedData map[string]interface{}
	inputByte, _ := json.Marshal(newStruct)
	json.Unmarshal(inputByte, &cleanedData)

	return cleanedData, nil
}

// Get the data without token
func GetTemporaryUserDataByAccessKey(c *fiber.Ctx) error {

	filter :=
		bson.M{"access_key": c.Params("access_key")}

	response, err := GetQueryResult("temporary_user", filter, int64(0), int64(2), nil)

	if err != nil {
		return BadRequest(err.Error())
	}

	return SuccessResponse(c, response)

}

func DeleteTempDocumentsBasedexpireTime() {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Define your filter for deletion
			currentTime := time.Now()
			filter := bson.M{
				"expire_on": bson.M{
					"$lt": time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC),
				},
			} // Delete records matching the filter
			_, err := database.GetConnection().Collection("temporary_user").DeleteMany(ctx, filter)
			if err != nil {
				fmt.Printf("Error deleting records: %v\n", err)
			}

			// Sleep for an interval
			time.Sleep(time.Hour) // Adjust the interval as needed
		}
	}
}

func MasterAggregationPipeline(request PaginationRequest, c *fiber.Ctx) []bson.M {
	Pipeline := []bson.M{}

	if len(request.Filter) > 0 {
		FilterConditions := BuildAggregationPipeline(request.Filter, "")
		Pipeline = append(Pipeline, FilterConditions)
	}

	if len(request.Sort) > 0 {
		sortConditions := buildSortConditions(request.Sort)
		Pipeline = append(Pipeline, sortConditions)
	}

	return Pipeline
}

func GroupDataBasedOnRules(c *fiber.Ctx) error {
	OrgId := c.Get("OrgId")
	if OrgId == "" {
		return BadRequest("Organization Id missing")
	}

	filter := bson.M{"group_name": c.Params("groupname")}

	// var response map[string]interface{}
	response, err := FindDocs(OrgId, "group", filter)
	if err != nil {
		return err
	}

	delete(response, "_id")
	delete(response, "group_name")
	delete(response, "grouptype")
	delete(response, "ref_collection")
	delete(response, "status")

	// fmt.Println(response["filter"])
	// responses := response["filter"].(Primitive.A)
	// result := PaginationRequest{
	// 	Filter: responses,
	// }
	var result PaginationRequest
	hello, _ := json.Marshal(response)

	if err := json.Unmarshal(hello, &result); err != nil {
		// fmt.Println("Error:", err)
		// return
	}

	return c.JSON(result)
}

// DatasetsConfig -- METHOD PURPOSE handle requests related to dataset configuration, including building aggregation pipelines
func DatasetsConfig(c *fiber.Ctx) error {
	//Get the OrgId from Header
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}
	//TO Bind the Value from Body
	var inputData DataSetConfiguration
	if err := c.BodyParser(&inputData); err != nil {
		if cmdErr, ok := err.(mongo.CommandError); ok {
			return BadRequest(cmdErr.Message)
		}

	}

	// BuildPipeline -- Create a Filter Pipeline from Body Content
	Data, PreviewResponse := BuildPipeline(orgId, inputData)
	// Set the DatasetName to _id for unique
	Data.Id = inputData.DataSetName
	// Params options -- options is  insert the data to Db
	// if options empty is preview the data
	if c.Params("options") == "Insert" {

		res, err := database.GetConnection().Collection("dataset_config").InsertOne(ctx, Data)
		if err != nil {
			return BadRequest("dataset name already used")
		}
		return SuccessResponse(c, res)

	}

	return SuccessResponse(c, PreviewResponse)

}

// BuildPipeline    -- METHOD PURPOSE  build a comprehensive MongoDB aggregation pipeline for querying and aggregating data
func BuildPipeline(orgId string, inputData DataSetConfiguration) (DataSetConfiguration, fiber.Map) {
	//append the pipelien from child Pipeline
	Pipeline := []bson.M{}
	//Every If condtion for if Data is here that time only that func work
	if len(inputData.DataSetBaseCollectionFilter) > 0 {
		Pipelines := BuildAggregationPipeline(inputData.DataSetBaseCollectionFilter, inputData.DataSetBaseCollection)
		Pipeline = append(Pipeline, Pipelines)

	}

	if len(inputData.DataSetJoinCollection) > 0 {
		lookupData := ExecuteLookupQueryData(inputData, inputData.DataSetBaseCollection)
		Pipeline = append(Pipeline, lookupData...)
	}

	if len(inputData.CustomColumn) > 0 {
		createCustomColumns := CreateCusotmColumns(Pipeline, inputData.CustomColumn, inputData.DataSetBaseCollection)
		Pipeline = append(Pipeline, createCustomColumns...)
	}

	if len(inputData.Aggregation) > 0 {
		AggregationData := buildDynamicAggregationPipeline(inputData.Aggregation)
		Pipeline = append(Pipeline, AggregationData...)
	}

	if len(inputData.Filter) > 0 {
		filterPipelines := BuildAggregationPipeline(inputData.Filter, inputData.DataSetBaseCollection)
		Pipeline = append(Pipeline, filterPipelines)

	}

	if len(inputData.SelectedList) > 0 {
		selectedColumns := CreateSelectedColumn(inputData.SelectedList, inputData.DataSetBaseCollection)
		Pipeline = append(Pipeline, selectedColumns...)
	}
	// filter pipeline convert the byte
	marshaldata, err := json.Marshal(Pipeline)
	if err != nil {
		return DataSetConfiguration{}, nil
	}
	// marshaldata  variable -- filter byte  convert the string
	pipelinestring := string(marshaldata)
	// set the inputData.Pipeline  -- store the data form converted string pipeine
	inputData.Pipeline = pipelinestring

	// Filter Params for to replace the string to convert to pipeline again
	if len(inputData.FilterParams) > 0 {
		inputData.Reference_pipeline = pipelinestring
		pipelinestring := createFilterParams(inputData.FilterParams, pipelinestring)
		// if filter params here that time to replace the old pipeline
		inputData.Pipeline = pipelinestring
		// convert the pipeline
		err := bson.UnmarshalExtJSON([]byte(pipelinestring), true, &Pipeline)
		if err != nil {
			fmt.Println("Error parsing pipeline:", err)

		}

	}

	finalpipeline := []bson.M{}
	Updatedpipeline := createQueryPipeline(Pipeline)
	//final pagination TO add the Filter
	Data, _ := Updatedpipeline.([]primitive.M)
	finalpipeline = append(finalpipeline, Data...)

	PagiantionPipeline := PagiantionPipeline(inputData.Start, inputData.End)
	finalpipeline = append(finalpipeline, PagiantionPipeline)

	// Get the Data form Db
	Response, err := GetAggregateQueryResult(inputData.DataSetBaseCollection, finalpipeline)
	if err != nil {
		return DataSetConfiguration{}, nil
	}
	// this PreviewResponse
	PreviewResponse := fiber.Map{
		"status":   "success",
		"data":     Response,
		"pipeline": Pipeline,
	}

	return inputData, PreviewResponse

}

// Insert the Data and return map
func InsertDataDb(orgId string, inputData interface{}, collectionName string) (fiber.Map, error) {

	res, err := database.GetConnection().Collection(collectionName).InsertOne(ctx, inputData)
	InsertResponse := fiber.Map{
		"status":  "success",
		"message": "Data Added Successfully",
		"data": fiber.Map{
			"InsertedID": res.InsertedID,
		},
	}
	return InsertResponse, err
}

// DatasetsRetrieve  -- METHOD PURPOSE Get the Filter pipeline in Db to show the data
func DatasetsRetrieve(c *fiber.Ctx) error {
	//OrgId oming from Header
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}
	//Params
	datasetname := c.Params("datasetname")

	filter := bson.M{"dataSetName": datasetname}
	// Find the Data from Db
	response, err := FindDocs(orgId, "dataset_config", filter)
	if err != nil {
		return BadRequest("Invalid  Params value")
	}

	var ResponseData map[string]interface{}
	//Marshal the data from find the Document
	marshaldata, err := json.Marshal(response)
	if err != nil {
		return BadRequest("Failed to Marshal ")
	}

	// Unmarshal -- after Unmarshal to map[string]interface{} convert
	if err := json.Unmarshal(marshaldata, &ResponseData); err != nil {
		return BadRequest("Invalid Body Content from MarshalData")

	}
	var requestBody PaginationRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return BadRequest("Invalid Body")
	}

	// Set the global variable for append the value from filter params
	var Dbpipelinestring string
	// if filterparams is here that time to replace the value get he ResponseData from referencepipeline
	if len(requestBody.FilterParams) > 0 {
		pipelinestring := createFilterParams(requestBody.FilterParams, ResponseData["referencepipeline"].(string))
		Dbpipelinestring = pipelinestring
	} else {
		// if not there return the pipeine storeed variable data return from ResponseData
		Dbpipelinestring = ResponseData["pipeline"].(string)
	}

	//Get the Collection Name in Database
	CollectionName := ResponseData["dataSetBaseCollection"].(string)
	//Body Filter storing to struct

	// Parse the provided string into a slice of BSON documents for the pipeline.
	pipeline := []primitive.M{}
	err = json.Unmarshal([]byte(Dbpipelinestring), &pipeline)
	if err != nil {
		return BadRequest("Cannot Find the String")

	}

	//finalpipeline -- Build the Final append filter pipeline
	var finalpipeline []bson.M
	//UpdateDatatypes -- To build the Pipeline from pipeline variable
	Updatedpipeline := UpdateDatatypes(pipeline)

	finalpipeline = append(finalpipeline, Updatedpipeline...)

	//Body Filter Pipeline making
	filterpipeline := MasterAggregationPipeline(requestBody, c)
	//To combine the pipeline filter and basefilter
	finalpipeline = append(finalpipeline, filterpipeline...)

	//final pagination TO add the Filter
	PagiantionPipeline := PagiantionPipeline(requestBody.Start, requestBody.End)
	finalpipeline = append(finalpipeline, PagiantionPipeline)

	// To Get the Data from Db
	Response, err := GetAggregateQueryResult(CollectionName, finalpipeline)
	if err != nil {
		InternalServerError(err.Error())
	}

	return SuccessResponse(c, Response)
}

// UpdateDatatypes    --METHOD  Get the match object and to build the mongo Query
func UpdateDatatypes(pipeline []bson.M) []bson.M {
	output := []bson.M{}
	for _, stage := range pipeline {
		if matchStage, ok := stage["$match"]; ok {
			// To Pass the interface{} to $match data for datatype convertion
			matchedPipeline := createQueryPipeline(matchStage)
			// to append the  convert datatype then add inside the match if $match is not there else work it..
			output = append(output, bson.M{"$match": matchedPipeline})
		} else {
			output = append(output, stage)
		}
	}

	return output
}

/*
	createQueryPipeline -- METHOD To change the value Datatype and return the pipeline format

Recusively call the  Method for Datatype converntiuon
*/

func createQueryPipeline(data interface{}) interface{} {
	// Check the Every DataType to incoming
	switch dataType := data.(type) {
	case map[string][]interface{}:
		var outputArray []interface{}
		for _, value := range dataType {
			for _, item := range value {
				outputArray = append(outputArray, createQueryPipeline(item))
			}
		}
		return outputArray
	case map[string]interface{}:
		valueMap := dataType
		for k := range valueMap {
			valueMap[k] = createQueryPipeline(valueMap[k])
		}
		return valueMap
	case []interface{}:
		var outputArray []interface{}
		for _, i := range dataType {
			outputArray = append(outputArray, createQueryPipeline(i))
		}
		return outputArray
	default:
		// if return final interface{} to ge the convert the data type to  ConvertToDataType
		return ConvertToDataType(data, reflect.TypeOf(data).String())
	}
}

// UpdateDataset  --METHOD Update the Dataset_config collection to store the data with pipeline
func UpdateDataset(c *fiber.Ctx) error {
	// Get the OrgId from header
	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}
	datasetname := c.Params("datasetname")
	//Params
	filter := bson.M{"dataSetName": datasetname}
	//Update body to bind the  DataSetConfiguration
	var inputData DataSetConfiguration
	if err := c.BodyParser(&inputData); err != nil {
		return BadRequest("Invalid Body Content")
	}
	// Global Variable set the For response
	var Response fiber.Map
	//Build the Pipeline
	Data, Response := BuildPipeline(orgId, inputData)
	// update the data set to Db
	Response, err := UpdateDataToDb(orgId, filter, Data, "dataset_config")
	if err != nil {
		return BadRequest("Failed to insert data into the database")

	}

	return SuccessResponse(c, Response)
}
