package profile

import "time"

// Profile data
type Profile struct {
	Login string
	Name  string
	Web   string
}

// FetchProfile Devuelve información de Usuario
func FetchProfile() *Profile {
	// Un delay para simular tiempo de espera en llamadas remotas
	time.Sleep(1 * time.Second)

	return &Profile{
		Login: "nmarsollier",
		Name:  "Nestor Marsollier",
		Web:   "https://github.com/nmarsollier/profile",
	}
}
