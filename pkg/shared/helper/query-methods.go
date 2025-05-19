package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"introme-api/pkg/shared/database"
)

var updateOpts = options.Update().SetUpsert(true)
var findUpdateOpts = options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
var ctx = context.Background()

func GetAggregateQueryResult(collectionName string, query interface{}) ([]bson.M, error) {
	response, err := ExecuteAggregateQuery(collectionName, query)
	if err != nil {
		return nil, err
	}
	var result []bson.M
	//var result map[string][]Config
	if err = response.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func ExecuteAggregateQuery(collectionName string, query interface{}) (*mongo.Cursor, error) {
	cur, err := database.GetConnection().Collection(collectionName).Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func GetQueryResult(collectionName string, query interface{}, page int64, limit int64, sort interface{}) ([]bson.M, error) {
	response, err := ExecuteQuery(collectionName, query, page, limit, sort)
	if err != nil {
		return nil, err
	}

	var result []bson.M
	//var result map[string][]Config
	if err = response.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetQueryInBetweenId(orgId string, collectionName string, options *options.FindOptions, startId string, endId string) ([]bson.M, error) {
	var result []bson.M
	// Construct a filter to find documents with IDs between abc123 and xyz789
	filter := bson.M{
		"_id": bson.M{
			"$gte": startId,
			"$lte": endId,
		},
	}
	cur, err := database.GetConnection().Collection(collectionName).Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func ExecuteHistoryInsertMany(orgId string, collectionName string, docs []interface{}) (*mongo.InsertManyResult, error) {
	result, err := database.GetConnection().Collection(collectionName+"_history").InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ExecuteDeleteManyByIds(orgId string, collectionName string, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := database.GetConnection().Collection(collectionName).DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ExecuteQuery(collectionName string, query interface{}, page int64, limit int64, sort interface{}) (*mongo.Cursor, error) {
	pageOptions := options.Find()
	pageOptions.SetSkip(page)   //0-i
	pageOptions.SetLimit(limit) // number of records to return
	if sort != nil {
		pageOptions.SetSort(sort)
	}
	response, err := database.GetConnection().Collection(collectionName).Find(ctx, query, pageOptions)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ExecuteFindAndModifyQuery(collectionName string, filter interface{}, data interface{}) (bson.M, error) {
	var result bson.M

	err := database.GetConnection().Collection(collectionName).FindOneAndUpdate(ctx, filter, data, findUpdateOpts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetReportQueryResult(orgId string, collectioinName string, req ReportRequest) ([]bson.M, error) {
	//build filter query
	query := make(map[string]interface{})
	//Check emp id
	if req.EmpId != "" {
		query["eid"] = req.EmpId
	}

	//check emp id
	if len(req.EmpIds) > 0 {
		query["eid"] = bson.M{"$in": req.EmpIds}
	}
	//if date filter presented or not
	if req.DateColumn == "" { // start & end filter
		if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
			query["start_date"] = bson.M{"$gte": req.StartDate, "$lte": req.EndDate}
			query["end_date"] = bson.M{"$gte": req.StartDate, "$lte": req.EndDate}
		} else if !req.StartDate.IsZero() && req.EndDate.IsZero() {
			query["start_date"] = bson.M{"$gte": req.StartDate}
		} else if req.StartDate.IsZero() && !req.EndDate.IsZero() {
			query["end_date"] = bson.M{"$lte": req.EndDate}
		}
	} else { // in between date filter
		if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
			query[req.DateColumn] = bson.M{"$gte": req.StartDate, "$lte": req.EndDate}
		} else if !req.StartDate.IsZero() && req.EndDate.IsZero() {
			query[req.DateColumn] = bson.M{"$gte": req.StartDate}
		} else if req.StartDate.IsZero() && !req.EndDate.IsZero() {
			query[req.DateColumn] = bson.M{"$lte": req.EndDate}
		}
	}
	if req.Type != "" {
		query["type"] = req.Type
	}
	if req.Status != "" {
		query["status"] = req.Status
	}
	return GetQueryResult(collectioinName, query, int64(0), int64(200), nil)
}

func getCondition(field string, value string) bson.M {
	condition := []string{"$" + field, value}
	return bson.M{
		"$sum": bson.M{
			"$cond": []interface{}{bson.M{"$eq": condition}, 1, 0},
		},
	}
}

func Updateformodel(c *fiber.Ctx) error {

	orgId := c.Get("OrgId")
	if orgId == "" {
		return BadRequest("Organization Id missing")
	}

	// fmt.Println("Insert The Collections")
	collectionName := c.Params("collectionName")

	// If the ID is not a valid ObjectID, search using the ID as a string
	filter := DocIdFilter(c.Params("id"))

	var inputData map[string]interface{}
	if err := c.BodyParser(&inputData); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:   fiber.StatusBadRequest,
			ErrorMsg: "Error parsing request body",
		})
	}
	update := bson.M{
		"$set": inputData,
	}
	// Update data in the collection
	res, err := database.GetConnection().Collection(collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		response := Response{
			Status:   fiber.StatusInternalServerError,
			ErrorMsg: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := fiber.Map{
		"status": 200,
		// "data":      []map[string]interface{}{inputData},
		"data":      res.UpsertedID,
		"error_msg": "",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func FindDocs(collection string, filter interface{}) (map[string]interface{}, error) {

	var result map[string]interface{}
	err := database.GetConnection().Collection(collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle case when no document matches the filter
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

func BuildAggregationPipeline(inputData []FilterCondition, BasecollectionName string) bson.M {
	var matchConditions []bson.M

	for _, filter := range inputData {
		for _, condition := range filter.Conditions {
			Finals := GenerateAggregationPipeline(condition, BasecollectionName)
			matchConditions = append(matchConditions, Finals...)
		}
	}

	var clause bson.M
	if len(matchConditions) > 0 {
		if inputData[0].Clause == "OR" {
			clause = bson.M{"$or": matchConditions}
		} else if inputData[0].Clause == "AND" {
			clause = bson.M{"$and": matchConditions}
		}
	}

	return bson.M{"$match": clause}
}

func GenerateAggregationPipeline(condition ConditionGroup, basecollection string) []bson.M {
	conditions := []bson.M{}

	//* If Nested Conditions Here that time Recursively load the filter
	if len(condition.Conditions) > 0 {
		nestedConditions := []bson.M{}
		for _, nestedCondition := range condition.Conditions {
			Nested := GenerateAggregationPipeline(nestedCondition, basecollection)

			nestedConditions = append(nestedConditions, Nested...)
		}

	}

	column := condition.Column
	value := condition.Value
	reference := condition.ParentCollectionName

	//if basecollection is empty we use directly use columnName
	if basecollection == "" {
		column = condition.Column
	} else if condition.ParentCollectionName == "" { //If ParentCollectioName is  empty we use directly use columnName
		column = condition.Column
	} else if basecollection != condition.ParentCollectionName { //If basecollection and ParentCollectionName is not equal that time suse refence variable for DOT
		column = reference + "." + fmt.Sprint(column)
	}

	//What are the Opertor is here mention that  map
	operatorMap := map[string]string{
		"EQUALS":             "$eq",
		"NOTEQUAL":           "$ne",
		"CONTAINS":           "$regex",
		"NOTCONTAINS":        "$regex",
		"STARTSWITH":         "$regex",
		"ENDSWITH":           "$regex",
		"LESSTHAN":           "$lt",
		"GREATERTHAN":        "$gt",
		"LESSTHANOREQUAL":    "$lte",
		"GREATERTHANOREQUAL": "$gte",
		"INRANGE":            "$gte",
		"BLANK":              "$exists",
		"NOTBLANK":           "$exists",
		"EXISTS":             "$exists",
		"IN":                 "$in",
	}

	//OpertorMap check we Sended In body to map
	if operator, exists := operatorMap[condition.Operator]; exists {
		// conditionValue := ConvertToDataType(value, condition.Type, valueType, pipelineBuild)

		conditionValue := ConvertToDataType(value, condition.Type)
		if condition.Operator == "INRANGE" || condition.Operator == "IN_BETWEEN" {
			if condition.Type == "date" || condition.Type == "time.Time" {
				dateValues, isDate := value.([]interface{})
				if isDate && len(dateValues) == 2 {
					startDateValue, startOK := dateValues[0].(string)
					endDateValue, endOK := dateValues[1].(string)
					if startOK && endOK {
						startDate, startErr := time.Parse(time.RFC3339, startDateValue)
						endDate, endErr := time.Parse(time.RFC3339, endDateValue)
						if startErr == nil && endErr == nil {
							startOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
							endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
							conditions = append(conditions, bson.M{column: bson.M{
								"$gte": startOfDay,
								"$lte": endOfDay,
							}})
						}
					}
				}
			} else {
				rangeValues, ok := value.([]interface{})
				if ok && len(rangeValues) == 2 {
					minValue := rangeValues[0]
					maxValue := rangeValues[1]
					conditions = append(conditions, bson.M{column: bson.M{"$gte": minValue, "$lte": maxValue}})
				}
			}
		}

		if condition.Operator == "BLANK" {
			conditions = append(conditions, bson.M{column: bson.M{operator: false}})
		} else if condition.Operator == "EXISTS" {
			conditions = append(conditions, bson.M{column: bson.M{operator: conditionValue}})
		} else if condition.Operator == "NOTCONTAINS" {
			pattern := fmt.Sprintf("^(?!.*%s)", condition.Value)
			// conditions = append(conditions, bson.M{condition.Column: bson.M{"$regex": pattern}})
			conditions = append(conditions, bson.M{condition.Column: bson.M{operator: pattern}})
		} else if condition.Operator == "NOTBLANK" {
			conditions = append(conditions, bson.M{column: bson.M{operator: true, "$ne": nil}})
		} else if condition.Operator == "IN" {

			conditions = append(conditions, bson.M{column: bson.M{operator: value.([]interface{})}})

		} else {
			conditions = append(conditions, bson.M{column: bson.M{operator: conditionValue}})
		}
	}

	//clause Binding
	if condition.Clause == "AND" {
		conditions = append(conditions, bson.M{"$and": conditions})
	} else if condition.Clause == "OR" {
		conditions = append(conditions, bson.M{"$or": conditions})
	}
	return conditions
}

// PagiantionPipeline -- METHOD Pagination return set of Limit data return
func PagiantionPipeline(start, end int) bson.M {
	// Get the Default value from env file
	startValue, _ := strconv.Atoi(os.Getenv("DEFAULT_START_VALUE"))
	endValue, _ := strconv.Atoi(os.Getenv("DEFAULT_LIMIT_VALUE"))

	//param is empty set the Default value

	if start == 0 {
		start = startValue
	}
	if end == 0 {
		end = endValue
	}

	// return the bson for pagination
	return bson.M{
		"$facet": bson.D{
			{"response",
				bson.A{
					bson.D{{"$skip", start}},
					bson.D{{"$limit", end}},
				},
			},
			{"pagination",
				bson.A{
					bson.D{{"$count", "totalDocs"}},
				},
			},
		},
	}
}

// ConvertToDataType --METHOD Build the Datatype from Paramters
func ConvertToDataType(value interface{}, DataType string) interface{} {
	// Check the data type and perform the corresponding conversion.
	if DataType == "time.Time" || DataType == "date" {
		// If the data type is time.Time, attempt to parse the value as a string in RFC3339 format.
		if valStr, ok := value.(string); ok {
			t, err := time.Parse(time.RFC3339, valStr)
			if err == nil {
				// If parsing is successful, return a time.Time value with truncated seconds.
				StartedDay := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
				return StartedDay
			}
		}
	} else if DataType == "string" || DataType == "text" {
		// If the data type is string or text, attempt to parse the value as a string.
		if valStr, ok := value.(string); ok {
			t, err := time.Parse(time.RFC3339, valStr)
			// If parsing as time is successful, return a time.Time value with truncated seconds.
			if err == nil {
				StartedDay := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
				return StartedDay
			} else {
				// If parsing as time fails, return the original string value.
				return valStr
			}
		}
	} else if DataType == "boolean" || DataType == "bool" {
		// If the data type is boolean or bool, attempt to cast the value to a boolean.
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	// If the data type is not recognized or conversion is not possible, return the original value.
	return value
}
func AddDurationToDate(durationStr string) (time.Time, error) {
	now := time.Now()
	durationStr = strings.ToLower(strings.TrimSpace(durationStr))

	// Helper functions to set the start or end of the desired periods
	startOfDay := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}
	endOfDay := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	}
	startOfMonth := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	}
	endOfMonth := func(t time.Time) time.Time {
		return startOfMonth(t.AddDate(0, 1, 0)).Add(-time.Nanosecond)
	}
	startOfYear := func(t time.Time) time.Time {
		return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	}
	endOfYear := func(t time.Time) time.Time {
		return time.Date(t.Year(), time.December, 31, 23, 59, 59, 999999999, t.Location())
	}
	startOfWeek := func(t time.Time) time.Time {
		offset := (int(t.Weekday()) + 6) % 7
		return t.AddDate(0, 0, -offset)
	}
	endOfWeek := func(t time.Time) time.Time {
		return startOfWeek(t).AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Split the duration string
	for _, part := range strings.Fields(durationStr) {
		if len(part) < 2 {
			return time.Time{}, fmt.Errorf("invalid format: %s", part)
		}

		sign := 1
		if part[0] == '-' {
			sign = -1
			part = part[1:]
		}

		// Handle special time settings based on suffix
		switch part {
		case "ds":
			now = startOfDay(now)
			continue
		case "de":
			now = endOfDay(now)
			continue
		case "ms":
			now = startOfMonth(now)
			continue
		case "me":
			now = endOfMonth(now)
			continue
		case "ys":
			now = startOfYear(now)
			continue
		case "ye":
			now = endOfYear(now)
			continue
		case "ws":
			now = startOfWeek(now)
			continue
		case "we":
			now = endOfWeek(now)
			continue
		}

		// Parse numeric duration and unit
		numPart := part[:len(part)-1]
		unit := part[len(part)-1:]

		value, err := strconv.Atoi(numPart)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration value: %s", numPart)
		}

		switch unit {
		case "d":
			now = now.AddDate(0, 0, value*sign)
		case "w":
			now = now.AddDate(0, 0, value*7*sign)
		case "m":
			now = now.AddDate(0, value*sign, 0)
		case "y":
			now = now.AddDate(value*sign, 0, 0)
		default:
			return time.Time{}, fmt.Errorf("unsupported unit: %s", unit)
		}
	}

	return now, nil
}

// UpdateDatasetConfig -- METHOD update the Data  to Db from filter and Data and collectionName from Param
func UpdateDataToDb(orgId string, filter interface{}, Data interface{}, collectionName string) (fiber.Map, error) {
	res, err := database.GetConnection().Collection(collectionName).UpdateOne(ctx, filter, Data)
	if err != nil {
		return nil, InternalServerError(err.Error())
	}
	UpdatetResponse := fiber.Map{
		"status":  "success",
		"message": "Update Successfully",
		"Data":    res.UpsertedID,
	}

	return UpdatetResponse, err
}
