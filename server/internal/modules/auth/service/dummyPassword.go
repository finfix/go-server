package service

// dummyPasswordSalt и dummyPasswordHash — заглушка для сравнения пароля несуществующего
// пользователя: проверка занимает столько же, сколько настоящая
var (
	dummyPasswordSalt = []byte("dummy-password-salt")
	dummyPasswordHash = mustDummyHash()
)
