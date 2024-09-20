package profile

type Profile interface {
	ID() string
	Login() string
	Name() string
	Web() string
}

// Profile data
type profile struct {
	id    string
	login string
	name  string
	web   string
}

func (p *profile) ID() string {
	return p.id
}

func (p *profile) Login() string {
	return p.login
}

func (p *profile) Name() string {
	return p.name
}

func (p *profile) Web() string {
	return p.web
}

func FetchProfile(id string) Profile {
	return &profile{
		id:    id,
		login: "nmarsollier",
		name:  "Nestor Marsollier",
		web:   "https://github.com/nmarsollier/profile",
	}
}
