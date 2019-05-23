package dtos

import (
	"github.com/thedevsaddam/govalidator"
)

type Pipeline struct {
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	OnError string           `json:"on_error"`
	Blocks  map[string]Block `blocks`
}

func (p *Pipeline) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"name":     []string{"required"},
		"type":     []string{"required", "in:plain"},
		"on_error": []string{"required", "in:stop,continue"},
	}
}

func (p *Pipeline) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"name":     []string{"Descriptive name for pipeline"},
		"type":     []string{"Type of pipeline: plain"},
		"on_error": []string{"Error handling"},
	}
}

type Block struct {
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

func (b *Block) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"type":   []string{"required"},
		"config": []string{"required"},
	}
}

func (b *Block) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"type":   []string{"namespaced block type"},
		"config": []string{"block specific config"},
	}
}
