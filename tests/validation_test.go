package tests

import (
	"ComicCollector/main/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestRequestBodyStrings(t *testing.T) {
	var requestBody struct {
		Name        string
		Description string
	}

	requestBody.Name = "test"
	requestBody.Description = "hehe :P"

	err := utils.ValidateRequestBody(requestBody)
	if err != nil {
		t.Error(err)
	}
}

func TestRequestBodyObjectIds(t *testing.T) {
	var requestBody struct {
		Name        string
		Description string
		Userid      primitive.ObjectID
		NilId       primitive.ObjectID
		IdArray     []primitive.ObjectID
	}

	// Test a valid body
	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Userid = primitive.NewObjectID()
	requestBody.NilId = primitive.NewObjectID()
	requestBody.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}

	err := utils.ValidateRequestBody(requestBody)
	if err != nil {
		t.Error(err)
	}

	// Test an invalid body
	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Userid = primitive.NewObjectID()
	requestBody.NilId = primitive.NilObjectID                                                  // This is invalid
	requestBody.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NilObjectID} // This is invalid

	err = utils.ValidateRequestBody(requestBody)
	if err == nil {
		t.Error(err)
	}

	// Test an invalid body
	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Userid = primitive.NewObjectID()
	requestBody.NilId = primitive.NewObjectID()
	requestBody.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NilObjectID} // This is invalid

	err = utils.ValidateRequestBody(requestBody)
	if err == nil {
		t.Error(err)
	}
}
