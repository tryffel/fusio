package dtos

import (
	"github.com/thedevsaddam/govalidator"
)

// Common format to include multiple devices / groups
type IdList struct {
	Ids []string `json:"ids"`
}

func (list *IdList) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"ids": []string{"required", "uuid_array"},
	}
}

func (list *IdList) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"ids": []string{"Array of valid ids to add"},
	}
}
