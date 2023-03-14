package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Message: "Hit the broker Bro",
	}

	_ =app.writeJSON(w,http.StatusOK,payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w,r,&requestPayload)

	if err != nil{
		app.errorJSON(w,err)
		return
	}

	switch requestPayload.Action{
	case "auth":
		app.authenticate(w,requestPayload.Auth)

	default:
		app.errorJSON(w,errors.New("unknown action"))
	}

}

func (app *Config) authenticate (w http.ResponseWriter, a AuthPayload){
	//Create json to send to auth micro service

	jsonData , _ := json.MarshalIndent(a,"","\t")
	 
	//call service

	request , err := http.NewRequest("POST","http://authentication-service/authenticate",bytes.NewBuffer(jsonData))

	if err != nil{
		app.errorJSON(w,err)
		return
	}

	client := &http.Client{}

	response , err := client.Do(request)

	if err != nil{
		app.errorJSON(w,err)
		return
	}

	defer response.Body.Close()
	//get back correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w,errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted{

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Print(bodyString, " .\n")

		
		app.errorJSON(w,errors.New("error calling auth service"))
		return
	}

	//create a variable to read response body into 

	var jsonFromService jsonResponse

	//decode json froms service

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil{
		app.errorJSON(w,err)
		return
	}

	if jsonFromService.Error{
		app.errorJSON(w,err,http.StatusUnauthorized)
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w , http.StatusAccepted , payload)
}