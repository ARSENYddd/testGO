package models

import (
	"database/sql"
	"errors"

	"example.com/auth_service/config"
)

var refreshTokens = map[string]string{}
var db *sql.DB

func SaveRefreshToken(userId, hashedToken string) error {
	refreshTokens[userId] = hashedToken
	return nil
}

func GetRefreshToken(userId string) (string, error) {
	if token, ok := refreshTokens[userId]; ok {
		return token, nil
	}
	return "", errors.New("refresh token not found")
}

func DeleteRefreshToken(userId string) error {
	if _, ok := refreshTokens[userId]; ok {
		delete(refreshTokens, userId)
		return nil
	}
	return errors.New("refresh token not found")
}
func GetUserIPFromDB(userId string) (string, error) {
	db, err := sql.Open("postgres", config.DBConnectionString)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var ip string
	err = db.QueryRow("SELECT ip_address FROM users WHERE user_id = $1", userId).Scan(&ip)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func UpdateRefreshToken(userId string, newRefreshTokenHash string) error {
	db, err := sql.Open("postgres", "user=test_db password=qwe dbname=qwe sslmode=disable host=localhost port=5432")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET refresh_token_hash = $1 WHERE id = $2", newRefreshTokenHash, userId)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE tokens SET refresh_token_hash = $1 WHERE id = $2", newRefreshTokenHash, userId)
	if err != nil {
		return err
	}

	return nil
}
