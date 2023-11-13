package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Id primitive.ObjectID `json:"_id" bson:"_id"` //data para recuperar el id automatico de mongo
	Task string `json:"task"`
	IsCompleted bool `default:"false"`
}

type UpdateTask struct {
	Task string `json:"task,omitempty" bson:"task,omitempty"`
	IsCompleted bool `default:"false,omitempty"`
}