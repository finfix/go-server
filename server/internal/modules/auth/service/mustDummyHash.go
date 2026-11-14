package service

import "pkg/passwordManager"

// mustDummyHash считает хэш заглушки один раз при старте
func mustDummyHash() []byte {
	hash, err := passwordManager.CreateNewPassword([]byte("dummy-password"), nil, dummyPasswordSalt)
	if err != nil {
		panic(err)
	}
	return hash
}
