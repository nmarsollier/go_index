package dialog

import (
	"encoding/json"
)

type dialogAction struct {
	Label  string `json:"label,omitempty"`
	Action string `json:"action,omitempty"`
}

type dialog struct {
	Title  *string       `json:"title,omitempty"`
	Accept *dialogAction `json:"accept,omitempty"`
}

type DialogBuilder struct {
	dialog
}

func NewBuilder() *DialogBuilder {
	return &DialogBuilder{dialog{}}
}

func (d *DialogBuilder) Title(value string) *DialogBuilder {
	return &DialogBuilder{dialog{
		Title:  &value,
		Accept: d.dialog.Accept,
	}}
}

func (d *DialogBuilder) AcceptAction(label, action string) *DialogBuilder {
	return &DialogBuilder{dialog{
		Title: d.dialog.Title,
		Accept: &dialogAction{
			Label:  label,
			Action: action,
		},
	}}
}

func (d *DialogBuilder) Build() string {
	result, _ := json.Marshal(d.dialog)
	return string(result)
}
