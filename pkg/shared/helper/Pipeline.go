package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// // Build match conditions for a filter

//	var TypeMaps = map[string]interface{}{
//		"string":    new(string),
//		"time":      new(time.Time),
//		"bool":      new(bool),
//		"int":       new(int),
//		"int32":     new(int32),
//		"int64":     new(int64),
//		"float64":   new(float64),
//		"time.Time": new(*time.Time),
//		"[]string":  new(string),
//		"[]time":    new(time.Time),
//		"[]bool":    new(bool),
//		"[]int":     new(int),
//		"[]int32":   new(int32),
//		"[]int64":   new(int64),
//		"[]float64": new(float64),
//	}

// var DataType = map[string]interface{}{
// 	"string":    new(string),
// 	"time":      new(time.Time),
// 	"bool":      new(bool),
// 	"int":       new(int),
// 	"int32":     new(int32),
// 	"int64":     new(int64),
// 	"float64":   new(float64),
// 	"time.Time": new(*time.Time),
// 	"[]string":  new([]string),
// 	"[]time":    new([]time.Time),
// 	"[]bool":    new([]bool),
// 	"[]int":     new([]int),
// 	"[]int32":   new([]int32),
// 	"[]int64":   new([]int64),
// 	"[]float64": new([]float64),
// }

// Build sort conditions for a filter
func buildSortConditions(sortCriteria []SortCriteria) bson.M {
	sortConditions := bson.M{}

	if len(sortCriteria) > 0 {
		sortCriteriaMap := bson.M{}
		for _, sort := range sortCriteria {
			order := 1 // Default to ascending order
			if sort.Sort == "desc" {
				order = -1 // Change to descending order
			}
			sortCriteriaMap[sort.ColID] = order
		}

		sortConditions["$sort"] = sortCriteriaMap
	}

	return sortConditions
}

func ExecuteLookupQueryData(Data DataSetConfiguration, basecollectionName string) []bson.M {
	var lookupDataPipeline []bson.M
	var previousAs string
	currentCollection := Data.DataSetBaseCollection

	for _, LookupData := range Data.DataSetJoinCollection {
		if LookupData.FromCollection == currentCollection {
			localField := LookupData.FromCollectionField

			if LookupData.FromCollection == previousAs {
				localField = previousAs + "." + LookupData.FromCollectionField
			}

			lookupStage := bson.M{
				"$lookup": bson.M{
					"from":         LookupData.ToCollection,
					"localField":   localField,
					"foreignField": LookupData.ToCollectionField,
					"as":           LookupData.ToCollection,
				},
			}
			lookupDataPipeline = append(lookupDataPipeline, lookupStage)

			unwindStage := bson.M{
				"$unwind": "$" + LookupData.ToCollection,
			}

			lookupDataPipeline = append(lookupDataPipeline, unwindStage)

			previousAs = LookupData.ToCollection
			currentCollection = LookupData.ToCollection

			if len(LookupData.Filter) > 0 {
				filterPipeline := BuildAggregationPipeline(LookupData.Filter, basecollectionName)
				lookupDataPipeline = append(lookupDataPipeline, filterPipeline)
			}
		}
	}

	return lookupDataPipeline
}

func CreateAggregationStage(column CustomColumn, Basecollection string) bson.M {
	var Pipeline bson.M
	// switch column.DataSetCustomAggregateFnName {
	if column.DataSetCustomAggregateFnName == "CONCAT" {
		Pipeline = addFieldsStage(column.DataSetCustomColumnName, bson.M{
			"$concat": generateConcat(column.DataSetCustomField, Basecollection)})

	} else if column.DataSetCustomAggregateFnName == "SUBTRACT" {
		Pipeline = addFieldsStage(column.DataSetCustomColumnName, bson.M{
			"$subtract": generateSub(column.DataSetCustomField, Basecollection),
		})
	} else if column.DataSetCustomAggregateFnName == "DIVIDE" {
		return addFieldsStage(column.DataSetCustomColumnName, bson.M{
			"$divide": generateSub(column.DataSetCustomField, Basecollection),
		})
	} else if column.DataSetCustomAggregateFnName == "MULTIPLY" {
		return addFieldsStage(column.DataSetCustomColumnName, bson.M{
			"$multiply": generateSub(column.DataSetCustomField, Basecollection),
		})
	} else if column.DataSetCustomAggregateFnName == "ADDITION" {
		return addFieldsStage(column.DataSetCustomColumnName, bson.M{
			"$add": generateSub(column.DataSetCustomField, Basecollection),
		})
	}

	return Pipeline
}

func buildDynamicAggregationPipeline(Aggregation []Aggregation) []bson.M {
	pipeline := []bson.M{}

	if len(Aggregation) == 0 {
		return pipeline
	}

	for _, aggregation_column := range Aggregation {
		GroupID := "$" + aggregation_column.AggGroupByField.Field
		group := bson.M{"_id": GroupID}
		AggColumnName := aggregation_column.AggColumnName
		FieldsName := "$" + aggregation_column.AggFieldName.Field
		AggFnName := aggregation_column.AggFnName

		if AggFnName == "SUM" {
			group[AggColumnName] = bson.M{"$sum": 1}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "MIN" {
			group[AggColumnName] = bson.M{"$min": FieldsName}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "MAX" {
			group[AggColumnName] = bson.M{"$max": FieldsName}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "PUSH" {
			group[AggColumnName] = bson.M{"$push": FieldsName}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "FIRST" {
			group[AggColumnName] = bson.M{"$first": FieldsName}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "LAST" {
			group[AggColumnName] = bson.M{"$last": FieldsName}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "COUNT" {
			group[AggColumnName] = bson.M{"$sum": 1}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		} else if AggFnName == "AVG" {
			group[AggColumnName] = bson.M{"$avg": FieldsName}
			group["doc"] = bson.M{"$first": "$$ROOT"}
		}

		pipeline = append(pipeline, bson.M{"$group": group})

		replaceRoot := bson.M{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"$mergeObjects": bson.A{
						bson.M{AggColumnName: "$" + AggColumnName},
						"$doc",
					},
				},
			},
		}
		pipeline = append(pipeline, replaceRoot)

		// return pipeline

	}
	return pipeline
}

func generateSub(fields []CustomField, Basecollection string) bson.A {
	expressions := bson.A{}

	for _, field := range fields {
		fieldName := field.FieldName
		fieldName = "$" + field.FieldName
		if Basecollection != field.ParentCollectionName {
			fieldName = "$" + field.ParentCollectionName + "." + field.FieldName
		}

		expressions = append(expressions, fieldName)
	}
	return expressions
}

func generateConcat(fields []CustomField, Basecollection string) bson.A {
	expressions := bson.A{}

	for i, field := range fields {

		FieldsName := field.FieldName
		FieldsName = "$" + field.FieldName
		if Basecollection != field.ParentCollectionName {
			FieldsName = "$" + field.ParentCollectionName + "." + field.FieldName
		}

		if i > 0 {
			expressions = append(expressions, " ")
		}
		expressions = append(expressions, FieldsName)
	}
	return expressions
}

func CreateCusotmColumns(Data []bson.M, CustomColumns []CustomColumn, Basecollection string) []bson.M {
	if len(CustomColumns) == 0 {
		return Data
	}

	aggregation := make([]bson.M, len(CustomColumns))
	// aggregation := make([]bson.M, 0)
	for i, column := range CustomColumns {
		aggregation[i] = CreateAggregationStage(column, Basecollection)
	}

	return aggregation
}

func CreateSelectedColumn(CustomColumns []SelectedListItem, BaseCollection string) []bson.M {
	fieldsToProject := bson.M{}

	for _, field := range CustomColumns {
		fieldName := field.Field

		fieldsToProject[fieldName] = 1
	}
	expressions := []bson.M{
		{
			"$project": fieldsToProject,
		},
	}

	return expressions
}

func addFieldsStage(dataSetCustomColumnName string, Fileds bson.M) bson.M {
	return bson.M{
		"$addFields": bson.M{
			dataSetCustomColumnName: Fileds,
		},
	}
}

func createFilterParams(FilterParams []FilterParam, Pipeline string) string {
	filterPipeline := Pipeline

	for _, Filter := range FilterParams {
		FindString := `{"parmasName":"` + Filter.ParamsName + `","parmsDataType":"` + Filter.ParamsDataType + `"}`

		var replacement string

		if Filter.DefaultValue != nil {
			c := reflect.TypeOf(Filter.DefaultValue).String()
			replacement = convertValueToDataType(c, Filter.DefaultValue)
		} else if Filter.Paramsvalue != nil {
			c := reflect.TypeOf(Filter.Paramsvalue).String()
			replacement = convertValueToDataType(c, Filter.Paramsvalue)

			// Handle unsupported data types
			if replacement == `"unsupported_data_type"` || replacement == `unsupported_data_type` {
				queryString, err := json.Marshal(Filter.Paramsvalue)
				if err != nil {
					fmt.Println("Error marshalling Paramsvalue:", err)
					continue
				}
				replacement = string(queryString)
			}
		}

		// Wrap in quotes if it's a string/date/time
		if Filter.ParamsDataType == "string" ||
			strings.ToLower(Filter.ParamsDataType) == "time.time" ||
			strings.ToLower(Filter.ParamsDataType) == "date" {

			// Only wrap if not already quoted
			if !strings.HasPrefix(replacement, `"`) {
				replacement = `"` + replacement + `"`
			}
		}

		// Replace in pipeline
		filterPipeline = strings.ReplaceAll(filterPipeline, FindString, replacement)
	}

	// fmt.Println(filterPipeline)
	return filterPipeline
}

// OLD
// func createFilterParams(FilterParams []FilterParam, Pipeline string) string {
// 	filterPipeline := Pipeline

// 	for _, Filter := range FilterParams {

// 		FindString := `{"ParamsName":"` + Filter.ParamsName + `","parmsDataType":"` + Filter.ParamsDataType + `"}`

// 		Replacement := `"` + Filter.DefaultValue + `"`
// 		filterPipeline = strings.ReplaceAll(filterPipeline, FindString, Replacement)

//		}
//		return filterPipeline
//	}
func convertValueToDataType(datatype string, defaultValue interface{}) string {
	var replaceValue string

	if datatype == "int" {
		if intValue, ok := defaultValue.(int); ok {
			replaceValue = strconv.Itoa(intValue)
		}
	} else if datatype == "int8" {
		if intValue, ok := defaultValue.(int8); ok {
			replaceValue = strconv.FormatInt(int64(intValue), 10)
		}
	} else if datatype == "int16" {
		if intValue, ok := defaultValue.(int16); ok {
			replaceValue = strconv.FormatInt(int64(intValue), 10)
		}
	} else if datatype == "int32" {
		if intValue, ok := defaultValue.(int32); ok {
			replaceValue = strconv.FormatInt(int64(intValue), 10)
		}
	} else if datatype == "int64" {
		if intValue, ok := defaultValue.(int64); ok {
			replaceValue = strconv.FormatInt(intValue, 10)
		}
	} else if datatype == "bool" {
		if boolValue, ok := defaultValue.(bool); ok {
			replaceValue = strconv.FormatBool(boolValue)
		}
	} else if datatype == "string" {
		if stringValue, ok := defaultValue.(string); ok {
			replaceValue = stringValue
		}
	} else if datatype == "float32" {
		if floatValue, ok := defaultValue.(float32); ok {
			replaceValue = strconv.FormatFloat(float64(floatValue), 'f', -1, 32)
		}
	} else if datatype == "float64" {
		if floatValue, ok := defaultValue.(float64); ok {
			replaceValue = strconv.FormatFloat(floatValue, 'f', -1, 64)
		}
	} else {
		replaceValue = "unsupported_data_type"
	}

	return replaceValue
}
