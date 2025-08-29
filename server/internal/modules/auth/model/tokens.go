package model

type Tokens struct {
	AccessToken  string `json:"accessToken"`  // Токен доступа
	RefreshToken string `json:"refreshToken"` // Токен восстановления доступа
}
