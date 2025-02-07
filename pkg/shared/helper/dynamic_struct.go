package helper

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"introme-api/pkg/shared/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

// TypeMap holds references to primitive types.
var TypeMap = map[string]interface{}{

	"string":    new(string),
	"time.Time": new(time.Time),
	"bool":      new(bool),
	"int":       new(int),
	"int32":     new(int32),
	"int64":     new(int64),
	"float64":   new(float64),
	// "time.Time": new(*time.Time),
	"[]string":  new([]string),
	"[]bool":    new([]bool),
	"[]int":     new([]int),
	"[]int32":   new([]int32),
	"[]int64":   new([]int64),
	"[]float64": new([]float64),
}

type Config struct {
	ModelName       string `json:"model_name" bson:"model_name"`
	ColumnName      string `json:"column_name" bson:"column_name"`
	Header          string `json:"header" bson:"header"`
	Type            string `json:"type" bson:"type"`
	Description     string `json:"description" bson:"description"`
	Tag             string `json:"tag" bson:"tag"`
	Status          string `json:"status" bson:"status"`
	Is_reference    bool   `json:"is_reference" bson:"is_reference"`
	Field           string `json:"field" bson:"field"`
	Collection_name string `json:"collection_name" bson:"collection_name"`
	Json_field      string `json:"json_field" bson:"json_field"`
}

type AggregateResult struct {
	ID     string   `bson:"_id"`
	Fields []Config `bson:"fields"`
}

func LoadDataModelFromDB(orgID string) map[string][]Config {
	var query = []bson.M{
		{
			"$match": bson.M{"status": "A"},
		},
		{
			"$group": bson.M{
				"_id": "$model_name",
				"fields": bson.M{
					"$push": bson.M{
						"column_name": "$column_name",
						"type":        "$type",
						"tag":         "$tag",
					},
				},
			},
		},
	}

	var result []AggregateResult
	cur, err := database.GetConnection().Collection("data_model").Aggregate(ctx, query)
	if err != nil {
		log.Errorf("Error aggregating data: %s", err.Error())
		return nil
	}
	defer cur.Close(ctx)

	if err = cur.All(ctx, &result); err != nil {
		log.Errorf("Error retrieving aggregation result: %s", err.Error())
		return nil
	}

	resultMap := make(map[string][]Config)
	for _, res := range result {
		resultMap[res.ID] = res.Fields
	}

	return resultMap
}

func loadModels(models map[string][]Config, key string) interface{} {
	// Check if the dynamic type for the model already exists in TypeMap.
	if model, exists := TypeMap[key]; exists {
		return model
	}

	// Create struct fields for the dynamic struct.
	var dynamicStruct []reflect.StructField
	for _, field := range models[key] {
		fieldName := field.ColumnName
		fieldType := TypeMap[field.Type]

		// Check if the type is already available in TypeMap.
		if _, exists := TypeMap[field.Type]; !exists {
			// Recursively load dependent models and their types.
			fieldType = loadModels(models, field.Type)
		}

		// Append the field definition to the dynamic struct.
		dynamicStruct = append(dynamicStruct,
			reflect.StructField{
				Name: fieldName,
				Type: reflect.TypeOf(fieldType),
				Tag:  reflect.StructTag(field.Tag),
			},
		)
	}

	// Create a struct type using the collected struct fields.
	obj := reflect.StructOf(dynamicStruct)
	objIns := reflect.New(obj).Interface()

	// Store the dynamically created struct type in TypeMap.
	TypeMap[key] = objIns

	// Return an instance of the dynamically created struct type.
	return objIns
}

func createDynamicTypes(data map[string][]Config) {
	for key := range data {
		// Load and create dynamic types for the models associated with the key.
		loadModels(data, key)
	}

}

func ServerInitstruct(orgIDs []string) {
	for _, orgID := range orgIDs {
		data := LoadDataModelFromDB(orgID)
		// Create dynamic schema types based on the data model configuration.
		createDynamicTypes(data)
	}
}

func CreateInstanceForCollection(collectionName string) (interface{}, map[string]string) {
	var validationErrors = make(map[string]string)

	exist := false
	for key := range TypeMap {
		if key == collectionName {
			exist = true
			break
		}
	}
	if !exist {
		validationErrors["error"] = fmt.Sprintf("INVALID COLLECTION NAME: %s", collectionName)
		return nil, validationErrors
	}

	objIns := reflect.New(reflect.TypeOf(TypeMap[collectionName])).Interface()

	return objIns, validationErrors
}

func Testing(c *fiber.Ctx) error {

	Response := DataChecking(c.Params("modelName"))

	Responses := fiber.Map{
		"Data": Response,
	}

	return c.JSON(Responses)
}

var DataTypes = map[string]string{
	"string":  "string",
	"bool":    "bool",
	"int":     "int",
	"int32":   "int32",
	"int64":   "int64",
	"float32": "float32",
	"float64": "float64",
}

func DataChecking(ModelName string) []bson.M {

	Data := DataTypeChecking(ModelName, "")

	ResponseData := NestedDatas(Data, ModelName, "")

	return ResponseData
}

// func NestedDatas(Data []bson.M, ModelName string) []bson.M {

// 	ResponseData := []bson.M{}
// 	var Flag bool

// 	for _, NestedData := range Data {

// 		Dataype := NestedData["type"].(string)

// 		if _, ok := DataTypes[Dataype]; !ok {
// 			Flag = true

// 			if Flag == true {
// 				MapData, _ := DataTypeChecking(Dataype, Dataype) //ModelName ,types   //Data address

// 				for _, inMap := range MapData {
// 					lowercaseString := strings.ToLower(inMap["name"].(string))

// 					fieldName := ModelName + "." + lowercaseString

// 					inMap["field_name"] = fieldName
// 					inMap["ParentCollectionName"] = ModelName
// 					NestedDataType := inMap["type"].(string)
// 					if _, ok := DataTypes[NestedDataType]; !ok {
// 						Flag = true
// 						NestedDatas(MapData, NestedDataType)

// 						// delete(nestedData...,NestedDataType)

// 					}
// 					ResponseData = append(ResponseData, inMap)

// 				}
// 			}

// 		} else {
// 			Flag = false
// 		}

// 	}

//		return ResponseData
//	}
func NestedDatas(Data []bson.M, ModelName string, typeToRecursivelyHandle string) []bson.M {
	ResponseData := []bson.M{}
	var Flag bool
	var typess string
	var fieldName string
	for _, NestedData := range Data {
		Dataype := NestedData["type"].(string)

		if _, ok := DataTypes[Dataype]; !ok {
			Flag = true

			if Flag {
				MapData := DataTypeChecking(Dataype, Dataype)

				for _, inMap := range MapData {
					fieldName = strings.ToLower(inMap["name"].(string))
					// fieldName = ModelName + "." + lowercaseString
					inMap["field_name"] = fieldName
					inMap["ParentCollectionName"] = ModelName
					typess = inMap["type"].(string)
					delete(NestedData, typess)
					delete(inMap, typess)
					ResponseData = append(ResponseData, inMap)

				}

				// Recursively process the specified type

				if _, ok := DataTypes[typess]; !ok {
					// fmt.Println(fieldName)

					recursivelyHandledData := DataTypeChecking(typess, typess)
					for _, recursivelyHandledDatass := range recursivelyHandledData {

						fieldName := fieldName + strings.ToLower(recursivelyHandledDatass["name"].(string))

						recursivelyHandledDatass["field_name"] = fieldName
						recursivelyHandledDatass["ParentCollectionName"] = ModelName
						delete(recursivelyHandledDatass, typess)
						//  det(recursivelyHandledDatass, "typess")
					}
					ResponseData = append(ResponseData, recursivelyHandledData...)

				}
			}
		} else {
			Flag = false
		}
	}

	return ResponseData
}

func DataTypeChecking(ModelName, Types string) []bson.M {
	// fmt.Println(ModelName)
	var pipeline = bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"status", "A"},
					{"model_name", ModelName},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$model_name"},
					{"Data",
						bson.D{
							{"$push",
								bson.D{
									{"name",
										bson.D{
											{"$toUpper",
												bson.D{
													{"$concat",
														bson.A{
															"$model_name",
															".",
															"$column_name",
														},
													},
												},
											},
										},
									},
									{"ParentCollectionName", "$model_name"},
									{"type", "$type"},
									{"field_name", "$json_field"},
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$unset", "_id"}},
		bson.D{{"$unwind", bson.D{{"path", "$Data"}}}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$Data"}}}},
	}

	if Types != "" {

		pipeline = append(pipeline, bson.D{{"$match", bson.D{{"type", bson.D{{"$ne", Types}}}}}})
	}

	Response, err := GetAggregateQueryResult("data_model", pipeline)

	if err != nil {
		// Handle the error
	}
	// DataChecking(Response)
	return Response
}
