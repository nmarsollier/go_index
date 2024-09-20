package profile

import (
	"fmt"
	"time"
)

// Profile data
type Profile struct {
	ID    string
	Login string
	Name  string
	Web   string
}

// FetchProfile Devuelve información de Usuario
func fetchProfile(id string) *Profile {
	fmt.Printf("Fetching Profile... %s \n", id)
	time.Sleep(1 * time.Second)
	return &Profile{
		ID:    id,
		Login: "nmarsollier",
		Name:  "Profile # " + id,
		Web:   "https://github.com/nmarsollier/profile",
	}
}
