package tests

import (
	"ComicCollector/main/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestRequestBodyStrings(t *testing.T) {
	var requestBody struct {
		Name        string `binding:"required"`
		Description string `binding:"required"`
	}

	requestBody.Name = "test"
	requestBody.Description = "hehe :P"

	err := utils.ValidateRequestBody(requestBody, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}

	requestBody.Name = ""
	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}

	var requestBody2 struct {
		Name        string
		Description string
	}

	requestBody2.Name = ""
	requestBody2.Description = ""

	err = utils.ValidateRequestBody(requestBody2, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}
}

func TestRequestBodyObjectIds(t *testing.T) {
	var requestBody struct {
		Name        string               `binding:"required"`
		Description string               `binding:"required"`
		Userid      primitive.ObjectID   `binding:"required"`
		NilId       primitive.ObjectID   `binding:"required"`
		IdArray     []primitive.ObjectID `binding:"required"`
	}

	// Test a valid body
	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Userid = primitive.NewObjectID()
	requestBody.NilId = primitive.NewObjectID()
	requestBody.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}

	err := utils.ValidateRequestBody(requestBody, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}

	// Test an invalid body
	requestBody.NilId = primitive.NilObjectID                                                  // This is invalid
	requestBody.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NilObjectID} // This is invalid

	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}

	// Test an invalid body
	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Userid = primitive.NewObjectID()
	requestBody.NilId = primitive.NewObjectID()
	requestBody.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NilObjectID} // This is invalid

	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}

	var requestBody2 struct {
		Name        string
		Description string
		Userid      primitive.ObjectID
		NilId       primitive.ObjectID
		IdArray     []primitive.ObjectID
	}

	requestBody2.Name = "test"
	requestBody2.Description = "hehe :P"
	requestBody2.Userid = primitive.NewObjectID()
	requestBody2.NilId = primitive.NewObjectID()
	requestBody2.IdArray = []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}

	err = utils.ValidateRequestBody(requestBody2, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}
}

func TestRequestBodyInts(t *testing.T) {
	var requestBody struct {
		Name        string `binding:"required"`
		Description string `binding:"required"`
		Age         int    `binding:"required"`
	}

	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Age = 20

	err := utils.ValidateRequestBody(requestBody, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}

	requestBody.Age = -1
	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}
}

func TestRequestBodyDateTime(t *testing.T) {
	var requestBody struct {
		Name        string             `binding:"required"`
		Description string             `binding:"required"`
		Age         int                `binding:"required"`
		BirthDate   primitive.DateTime `binding:"required"`
	}

	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Age = 20
	requestBody.BirthDate = primitive.NewDateTimeFromTime(time.Now())

	err := utils.ValidateRequestBody(requestBody, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}

	requestBody.BirthDate = primitive.DateTime(0)
	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}

	//var requestBody2 struct {
	//	Name        string
	//	Description string
	//	Age         int
	//	BirthDate   primitive.DateTime
	//}
	//
	//requestBody2.Name = "test"
	//requestBody2.Description = "hehe :P"
	//requestBody2.Age = 20
	//requestBody2.BirthDate = primitive.DateTime(0)
	//
	//err = utils.ValidateRequestBody(requestBody2)
	//if err != nil {
	//	t.Error("Expected no error, but got:", err)
	//}
}

func TestRequestBodyBytes(t *testing.T) {
	var requestBody struct {
		Name        string
		Description string
		Age         int
		Bytes       []byte
	}

	requestBody.Name = "test"
	requestBody.Description = "hehe :P"
	requestBody.Age = 20
	requestBody.Bytes = []byte("test")

	err := utils.ValidateRequestBody(requestBody, false)
	if err != nil {
		t.Error("Expected no error, but got:", err)
	}

	requestBody.Bytes = nil
	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}

	requestBody.Bytes = []byte("")
	err = utils.ValidateRequestBody(requestBody, false)
	if err == nil {
		t.Error(err)
	}
}
