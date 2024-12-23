package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}

}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}

	// does user exist
	userExists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("error checking if user exists:  %v", err)
	}
	if userExists {
		return fmt.Errorf("a user with the same username already exists")
	}

	// we know that the user does not exist
	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("error registering user: %v", err)
	}

	return nil
}
