package helper

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// !pending
func InitCustomValidator() {
	validate.RegisterValidation("within_duration", validateWithinDuration)
	validate.RegisterValidation("between_age", ValidateDateRange)
	validate.RegisterValidation("regexp", ValidateByCustomPattern)
	validate.RegisterValidation("within_durations", validateWithinDurations)

}

// ValidateDateRange is a custom validation function for checking if a person's age falls within a specified range
func ValidateDateRange(fl validator.FieldLevel) bool {
	// Convert the field value to a time.Time type
	dob, ok := fl.Field().Interface().(time.Time)
	if !ok {

		return false
	}

	tag := fl.Param()
	Tag := strings.Split(tag, "-")

	if len(Tag) != 2 {
		fmt.Printf("Invalid tag format\n")
		return false
	}

	minAgeStr := strings.TrimSuffix(Tag[0], "y")
	maxAgeStr := strings.TrimSuffix(Tag[1], "y")

	minAge, err := strconv.Atoi(minAgeStr)
	if err != nil {
		fmt.Printf("Failed to convert minAge: %v\n", err)
		return false
	}

	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		fmt.Printf("Failed to convert maxAge: %v\n", err)
		return false
	}

	// Calculate the age in years
	age := time.Since(dob).Hours() / 24 / 365.25

	// Check if the age is within the specified range
	if age < float64(minAge) || age > float64(maxAge) {
		fmt.Printf("Age is not within the valid range\n")
		return false
	}

	return true
}

func validateWithinDuration(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	if date, ok := value.(time.Time); ok {
		tagValue := fl.Param()

		re := regexp.MustCompile(`^(\d+)([mdyY])$`)
		match := re.FindStringSubmatch(tagValue)
		if len(match) != 3 {
			// Handle invalid tag value error
			return false
		}

		num, err := strconv.Atoi(match[1])
		if err != nil {
			// Handle number conversion error
			return false
		}
		unit := strings.ToLower(match[2])

		referenceDate := time.Now()

		// Calculate the earliest allowed date
		var earliestDate time.Time
		if unit == "m" {
			earliestDate = referenceDate.AddDate(0, -num, 0).Truncate(24 * time.Hour)
		} else if unit == "d" {
			earliestDate = referenceDate.AddDate(0, 0, -num).Truncate(24 * time.Hour)
		} else if unit == "y" {
			earliestDate = referenceDate.AddDate(-num, 0, 0).Truncate(24 * time.Hour)
		} else {
			// Handle invalid unit error
			return false
		}

		// Check if the date is within the allowed range
		if !date.Equal(earliestDate) && !(date.After(earliestDate) && date.Before(referenceDate)) {
			// Handle date not within range error
			return false
		}

		return true
	}

	// Handle invalid input type error
	return false
}

func ValidateByCustomPattern(fl validator.FieldLevel) bool {

	value := fl.Field().String()
	pattern := fl.Param()
	// Compile the provided regular expression pattern
	re := regexp.MustCompile(pattern)

	// Check if the value matches the regular expression
	return re.MatchString(value)
}

func validateWithinDurations(fl validator.FieldLevel) bool {
	// durationStr := fl.Field().String()

	// Parse the input duration string (e.g., "6m" for 6 months)
	// duration, err := time.ParseDuration(durationStr)
	// if err != nil {
	//     return false // Invalid duration format
	// }

	// // Get the current time
	// currentTime := time.Now()
	//
	// Calculate the time in the past based on the input duration
	// pastTime := currentTime.Add(-duration)

	// Retrieve the unit (e.g., "m" for month, "w" for week, etc.) from the duration string
	// unit := durationStr[len(durationStr)-1]

	// Calculate the start and end dates based on the unit
	// var start, end time.Time

	// switch unit {
	// case 'm':
	// For months, calculate the start and end of the current month
	//     start = time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, currentTime.Location())
	//     end = start.AddDate(0, 1, -1)
	// case 'w':
	//     // For weeks, calculate the start and end of the current week
	//     start = currentTime.AddDate(0, 0, -int(currentTime.Weekday()))
	//     end = start.AddDate(0, 0, 6)
	// case 'd':
	//     // For days, calculate the start and end of the current day
	//     start = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	//     end = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 999999999, currentTime.Location())
	// case 'y':
	//     // For years, calculate the start and end of the current year
	//     start = time.Date(currentTime.Year(), time.January, 1, 0, 0, 0, 0, currentTime.Location())
	//     end = time.Date(currentTime.Year(), time.December, 31, 23, 59, 59, 999999999, currentTime.Location())
	// default:
	//     return false // Invalid duration unit
	// }

	// Check if pastTime is within the specified range
	// return pastTime.After(start) && pastTime.Before(end)
	value := fl.Field().Interface()
	if date, ok := value.(time.Time); ok {
		tagValue := fl.Param()

		re := regexp.MustCompile(`^(\d+)([mdyY])$`)
		match := re.FindStringSubmatch(tagValue)
		if len(match) != 3 {
			// Handle invalid tag value error
			return false
		}

		num, err := strconv.Atoi(match[1])
		if err != nil {
			// Handle number conversion error
			return false
		}
		unit := strings.ToLower(match[2])

		referenceDate := time.Now()

		// Calculate the earliest allowed date
		var earliestDate time.Time
		if unit == "m" {
			earliestDate = referenceDate.AddDate(0, -num, 0).Truncate(24 * time.Hour)
		} else if unit == "d" {
			earliestDate = referenceDate.AddDate(0, 0, -num).Truncate(24 * time.Hour)
		} else if unit == "y" {
			earliestDate = referenceDate.AddDate(-num, 0, 0).Truncate(24 * time.Hour)
		} else {

			return false
		}
		if !date.Equal(earliestDate) && date.Before(earliestDate) {
			return false
		}

		// return earliestDate.After(num) && pastTime.Before(end)
		return true
	}

	// Handle invalid input type error
	return false
}

func ValidateStruct(inputStruct interface{}) map[string]string {
	// Create a map to store validation errors
	validationErrors := make(map[string]string)

	// Perform validation using the validator library
	err := validate.Struct(inputStruct)

	if err != nil {
		// Check if the error is related to validation setup or struct issues
		if _, ok := err.(*validator.InvalidValidationError); ok {
			validationErrors["error"] = "Validation setup or struct issues"
		} else {
			// Iterate through each validation error and extract relevant information
			for _, validationErr := range err.(validator.ValidationErrors) {
				// Get the name of the struct field where the validation error occurred
				fieldName := validationErr.StructField()

				// Get the specific validation tag that failed
				tagName := validationErr.Tag()

				// Get the value associated with the validation tag (if applicable)
				paramValue := validationErr.Param()

				// Create a user-friendly error message using the extracted information
				errorMessage := fmt.Sprintf(
					"Validation error for field %s: %s with param %s",
					fieldName, tagName, paramValue,
				)

				// Store the error message in the validationErrors map with the field name as the key
				validationErrors[fieldName] = errorMessage
			}
		}
	}

	// Return the map containing validation errors (if any)
	return validationErrors
}

/* FilterStructFieldsByJSON filters the fields of a struct based on the presence of field names
in a JSON field map. It also includes fields of nested structs that are not entirely empty.
The function returns a slice of filtered struct fields.
*/

func FilterStructFieldsByJSON(rv reflect.Value, jsonFieldMap map[string]interface{}) []reflect.StructField {
	// Initialize slices to store the filtered fields, nested struct fields, and checked variable names.
	var filteredFields []reflect.StructField
	var nestedStructFields []reflect.StructField
	var checkedVariableNames []string

	// Iterate over the fields of the input struct.
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		fieldValue := rv.Field(i)

		// Get the lowercase JSON field name from the tag.
		// fieldNameLower := strings.ToLower(string(field.Tag.Get("json")))

		fieldNameLower := string(field.Tag.Get("json"))
		// Check if the field name exists in the JSON field map.
		if _, exists := jsonFieldMap[fieldNameLower]; exists {
			// Handle fields that are pointers to nested structs.
			if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				newFields, _, variableName := ExtractNonEmptyFields(fieldValue.Elem(), field)

				// Keep track of checked variable names to avoid duplicate processing.
				checkedVariableNames = append(checkedVariableNames, variableName)
				nestedStructFields = append(nestedStructFields, newFields)
			}

			// Include non-nil fields in the filtered list.
			if !fieldValue.IsNil() {
				filteredFields = append(filteredFields, field)
			}
		}
	}

	// Iterate over the filtered fields and remove duplicates based on nested struct fields.
	var uniqueFilteredFields []reflect.StructField
	for _, field := range filteredFields {
		shouldKeep := true
		for _, varName := range checkedVariableNames {
			if field.Name == varName {
				shouldKeep = false
				break
			}
		}
		if shouldKeep {
			uniqueFilteredFields = append(uniqueFilteredFields, field)
		}
	}

	// Append the fields of nested structs to the filtered list.
	uniqueFilteredFields = append(uniqueFilteredFields, nestedStructFields...)

	// Return the final list of filtered struct fields.
	return uniqueFilteredFields
}

// CollectNonEmptyFields extracts non-empty fields from the given struct and creates a new struct type
// containing only those non-empty fields.
func ExtractNonEmptyFields(rv reflect.Value, fields reflect.StructField) (reflect.StructField, reflect.Value, string) {
	// Initialize a list to store non-empty field definitions.
	var nonEmptyFields []reflect.StructField

	// Create a new struct to store non-empty field values.
	nonEmptyValues := reflect.New(rv.Type()).Elem()

	// Iterate over the fields of the input struct.
	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		field := rv.Type().Field(i)

		// Check if the field is empty (nil pointer or empty struct).
		if (fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil()) ||
			(fieldValue.Kind() == reflect.Struct && fieldValue.NumField() == 0) {
			// Skip empty fields.
			continue
		}

		// Store the non-empty field definition.
		nonEmptyFields = append(nonEmptyFields,
			reflect.StructField{
				Name: field.Name,
				Type: field.Type,
				Tag:  field.Tag,
			})

		// Store the non-empty field value.
		nonEmptyValues.Field(i).Set(fieldValue)
	}

	// Create a new struct type containing only non-empty fields.
	dynamicStructType := reflect.StructOf(nonEmptyFields)

	// Update the field definitions with the new struct type.
	fields.Type = dynamicStructType

	// Store the name of the original field.
	FieldName := fields.Name

	// Return the updated field definitions, non-empty field values, and the original field name.
	return fields, nonEmptyValues, FieldName
}

// Insert
func vertifyInputStruct(rv reflect.Value, inputMap map[string]interface{}, errMap map[string]string) error {
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		// fieldTag := strings.ToLower(string(field.Tag.Get("json")))
		fieldTag := string(field.Tag.Get("json"))
		fieldValue := rv.Field(i)

		// Check if the JSON tag includes "omitempty".
		omitempty := false
		if strings.Contains(string(field.Tag.Get("json")), "omitempty") {
			omitempty = true
		}

		// If the field is not present in the input map and omitempty is set, skip it.
		if _, exists := inputMap[fieldTag]; !exists && omitempty {
			continue
		}

		if _, exists := inputMap[fieldTag]; !exists {
			errMap[fieldTag] = "Field missing"
		}

		if field.Type.Kind() == reflect.Struct {
			if err := vertifyInputStruct(rv.Field(i), inputMap, errMap); err != nil {
				return err
			}
			continue
		}

		//nested
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			for key, value := range inputMap {
				if key == fieldTag {
					if kk, isMap := value.(map[string]interface{}); isMap {
						if err := vertifyInputStruct(fieldValue.Elem(), kk, errMap); err != nil {
							return err
						}
					}
				}
			}
			continue
		}

		for key := range inputMap {
			keyExists := false
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Type().Field(i)
				// fmt.Println(field.Tag.Get("json"))
				// fieldTag := strings.ToLower(string(field.Tag.Get("json")))
				// fmt.Println(fieldTag)
				fieldTag := string(field.Tag.Get("json"))
				if fieldTag == key {
					keyExists = true
					break
				}
			}
			if !keyExists {
				errMap[key] = "Extra field"
			}
		}
	}

	if len(errMap) > 0 {
		return fmt.Errorf("Missing fields")
	}

	return nil
}

// todo checking
func updateFieldsWithParentKey(data map[string]interface{}, parentKey string, updatedData map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		currentKey := key
		if parentKey != "" {
			currentKey = parentKey + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			updateFieldsWithParentKey(v, currentKey, updatedData)
		case []interface{}:
			for index, item := range v {
				if nestedMap, isNestedMap := item.(map[string]interface{}); isNestedMap {
					nestedParentKey := fmt.Sprintf("%s[%d]", currentKey, index)
					updateFieldsWithParentKey(nestedMap, nestedParentKey, updatedData)
				}
			}
		default:
			updatedData[currentKey] = value
		}
	}
	// fmt.Println(updatedData)

	return updatedData
}
