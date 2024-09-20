package profile

import "fmt"

// Profile data
type Profile struct {
	ID    string
	Login string
	Name  string
	Web   string
}

// GetImage return the profile image resource
var GetImage = map[string]func(id string) string{
	"mobile": func(id string) string {
		return fmt.Sprintf("resource_%s", id)
	},
	"web": func(id string) string {
		return fmt.Sprintf("http://www.images.blablabla/images/%s", id)
	},
}

// IsValidDevice checks if device is valid
func IsValidDevice(device string) bool {
	_, ok := GetImage[device]
	return ok
}

// FetchProfile Devuelve información de Usuario
func FetchProfile(id string) *Profile {
	return &Profile{
		ID:    id,
		Login: "nmarsollier",
		Name:  "Nestor Marsollier",
		Web:   "https://github.com/nmarsollier/profile",
	}
}
