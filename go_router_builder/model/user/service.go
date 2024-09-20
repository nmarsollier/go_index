package user

import "time"

// User data
type User struct {
	Login  string
	Access string
}

// FetchUser Devuelve información de Usuario
func FetchUser() *User {
	// Un delay para simular tiempo de espera en llamadas remotas
	time.Sleep(1 * time.Second)

	return &User{
		Login:  "nmarsollier",
		Access: "ADMIN,USER",
	}
}
