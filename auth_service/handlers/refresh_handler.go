package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/auth_service/models"
	"example.com/auth_service/utils"
)

type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"user_id"`
}

func GetTokensHandler(writer http.ResponseWriter, request *http.Request) {
	var req RefreshTokensRequest

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" || req.UserId == "" {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	clientIP := request.RemoteAddr

	newAccessToken, newIP, err := utils.RefreshAccessToken(req.RefreshToken, req.UserId, clientIP)
	if err != nil {
		http.Error(writer, "Invalid refresh token33", http.StatusUnauthorized)
		return
	}
	oldIP, err := models.GetUserIPFromDB(req.UserId)
	if err != nil {
		return
	}
	if oldIP != newIP {
		go func() {
			err := utils.SendEmailWarning(req.UserId, oldIP, newIP)
			if err != nil {
				log.Printf("Failed to send email warning: %v", err)
			}
		}()
	}
	writer.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"access_token": newAccessToken,
	}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
