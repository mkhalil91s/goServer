package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct 
{
	Error bool `json error`
	Message string `json string`
	Data any `json:"data.omitEmpty"`
}
func (app *Config) Broker (w http.ResponseWriter, r *http.Request){

	payload := jsonResponse{
		Error : false,
		Message:  "Hi the broker",
	}

	out, _ := json.MarshalIndent(payload,"","\t")
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}