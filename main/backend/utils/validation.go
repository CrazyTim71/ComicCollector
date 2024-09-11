package utils

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/JoiHelper"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
)

func ContainsNilObjectID(array []primitive.ObjectID) bool {
	for _, id := range array {
		if id.IsZero() || id == primitive.NilObjectID {
			return true
		}
	}
	return false
}

func ValidateRequestBody(requestBody interface{}, checkObjectIdsForExistence bool) error {
	if requestBody == nil {
		return errors.New("request body cannot be nil")
	}

	// Iterate over all fields and validate them
	v := reflect.ValueOf(requestBody)
	t := reflect.TypeOf(requestBody)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		fieldTag := t.Field(i).Tag.Get("binding")

		// Check if the field is required
		isRequired := fieldTag == "required"

		switch fieldType.Kind() {
		case reflect.String:
			if field.String() == "" {
				if isRequired {
					return errors.New(t.Field(i).Name + " is required")
				} else {
					// Skip validation if the field is not required
					continue
				}
			}

			err := JoiHelper.UserInput.Validate(field.String())
			if err != nil {
				log.Println(err)
				return errors.New(v.Type().Field(i).Name + " is invalid")
			}
		case reflect.Int:
			if field.Int() < 0 {
				return errors.New(v.Type().Field(i).Name + " is invalid")
			}

		case reflect.Array: // primitive.ObjectID is an array
			if fieldType != reflect.TypeOf(primitive.ObjectID{}) {
				return errors.New(v.Type().Field(i).Name + " is an array of unknown type")
			}
			if field.Interface().(primitive.ObjectID).IsZero() {
				return errors.New(v.Type().Field(i).Name + " contains an invalid ObjectID")
			}
			if checkObjectIdsForExistence {
				// check if the objectID is valid
				var modelType string
				switch v.Type().Field(i).Name {
				case "BookType":
					modelType = "BookType"
				case "BookEdition":
					modelType = "BookEdition"
				case "Owners":
					modelType = "Owner"
				case "Locations":
					modelType = "Location"
				case "Publishers":
					modelType = "Publisher"
				case "Authors":
					modelType = "Author"
				default:
					return errors.New("unknown model type")
				}
				fmt.Print(modelType)

				if !operations.CheckIfExists(database.MongoDB, modelType, bson.M{"_id": field.Interface().(primitive.ObjectID)}) {
					return errors.New(v.Type().Field(i).Name + " does not exist")
				}
			}

		case reflect.Slice:
			if fieldType.Elem() == reflect.TypeOf(primitive.ObjectID{}) {
				if len(field.Interface().([]primitive.ObjectID)) == 0 {
					// only check if the field is required in case it is empty
					if isRequired {
						return errors.New(v.Type().Field(i).Name + " is empty")
					} // Skip validation if the field is not required
					continue
				}
				if fieldType.Elem() != reflect.TypeOf(primitive.ObjectID{}) {
					return errors.New(v.Type().Field(i).Name + " is a slice of unknown type")
				}
				if ContainsNilObjectID(field.Interface().([]primitive.ObjectID)) {
					return errors.New(v.Type().Field(i).Name + " contains an invalid ObjectID")
				}
			} else if field.Type().Elem().Kind() == reflect.Uint8 { // Check for []byte (slice of bytes)
				if len(field.Interface().([]byte)) == 0 {
					return errors.New(v.Type().Field(i).Name + " is empty")
				}
			} else {
				return errors.New(v.Type().Field(i).Name + " is a slice of unknown type")
			}
		case reflect.TypeOf(primitive.DateTime(0)).Kind():
			if field.Interface().(primitive.DateTime) == primitive.DateTime(0) {
				return errors.New(v.Type().Field(i).Name + " is invalid")
			}
		// TODO: add case for byte
		default:
			log.Println(fieldType.Kind())
			log.Println(reflect.TypeOf([]byte("")).Kind())
			panic("unhandled default case")
		}
	}

	return nil
}

func CleanEmptyFields(data interface{}) bson.M {
	result := bson.M{}
	v := reflect.ValueOf(data).Elem()

	// Loop through each field in the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Tag.Get("json") // Get the JSON tag name
		fieldType := field.Type()

		// If there's no json tag, or it's marked to be omitted, skip it
		if fieldName == "" || fieldName == "-" {
			continue
		}

		// Check if the field is a slice and empty. We only add non-empty slices to the result.
		if fieldType.Kind() == reflect.Slice {
			if len(field.Interface().([]primitive.ObjectID)) != 0 {
				result[fieldName] = field.Interface()
			}
			continue
		}

		if fieldType.Kind() == reflect.TypeOf(primitive.DateTime(0)).Kind() {
			if field.Interface().(primitive.DateTime) != primitive.DateTime(0) {
				result[fieldName] = field.Interface()
			}
			continue
		}

		if fieldType.Elem().Kind() == reflect.Uint8 { // Check for []byte
			if len(field.Interface().([]byte)) != 0 {
				result[fieldName] = field.Interface()
			}
			continue
		}

		// Check if the field is empty. We only add non-empty fields to the result.
		if !isEmptyValue(field) {
			result[fieldName] = field.Interface()
		}
	}

	return result
}

func isEmptyValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
