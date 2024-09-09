package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/auth_service/utils"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshTokenHandler(writer http.ResponseWriter, request *http.Request) {
	var tokenRequest RefreshTokenRequest

	err := json.NewDecoder(request.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	clientIP := request.RemoteAddr
	accessToken, refreshToken, err := utils.GenerateTokens(tokenRequest.UserID, clientIP)
	if err != nil {
		http.Error(writer, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
