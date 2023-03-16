package controller

import (
	"Explore1/model"
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, r *http.Request, req interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func SendErrorResponse(w http.ResponseWriter, message string) {
	var response model.ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendUnauthorizedResponse(w http.ResponseWriter) {
	var response model.ErrorResponse
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.WriteHeader(response.Status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
