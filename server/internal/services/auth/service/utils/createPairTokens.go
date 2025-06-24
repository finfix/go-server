package utils

import (
	"pkg/jwtManager"
	"server/internal/utils/auth"

	authModel "server/internal/services/auth/model"
)

func CreatePairTokens(userID uint32, deviceID string) (tokens authModel.Tokens, err error) {

	claims := auth.Claims{
		UserID:   userID,
		DeviceID: deviceID,
	}

	// Создаем Access token
	tokens.AccessToken, err = jwtManager.GenerateToken(jwtManager.AccessToken, claims)
	if err != nil {
		return tokens, err
	}

	// Создаем refresh token
	tokens.RefreshToken, err = jwtManager.GenerateToken(jwtManager.RefreshToken, claims)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
