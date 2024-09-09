package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"example.com/auth_service/config"
	"example.com/auth_service/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserId string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.StandardClaims
}

func GenerateTokens(userId, ip string) (string, string, error) {

	accessToken, err := generateAccessToken(userId, ip)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = models.UpdateRefreshToken(userId, refreshToken)
	if err != nil {
		return "", "", err
	}

	err = hashAndSaveRefreshToken(userId, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateAccessToken(userId, ip string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserId: userId,
		IP:     ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(config.JWTSecretKey)
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func hashAndSaveRefreshToken(userId, refreshToken string) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return models.SaveRefreshToken(userId, string(hashedToken))
}

func RefreshAccessToken(refreshToken, userId, ip string) (string, string, error) {

	hashedToken, err := models.GetRefreshToken(userId)
	log.Printf("Received refresh token: %s", refreshToken)
	log.Printf("User ID: %s", userId)
	log.Printf("Client IP: %s", ip)
	if err != nil {
		return "", "", errors.New("invalid refresh token1")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refreshToken))
	if err != nil {
		return "", "", errors.New("invalid refresh token2")
	}

	err = models.DeleteRefreshToken(userId)
	if err != nil {
		return "", "", errors.New("could not delete refresh token")
	}
	newAccessToken, err := generateAccessToken(userId, ip)
	if err != nil {
		return "", "", err
	}
	log.Printf("newAccessToken: %s", newAccessToken)
	oldIP := ip
	return newAccessToken, oldIP, nil
}
