package dtos

import (
	"github.com/thedevsaddam/govalidator"
	"github.com/tryffel/fusio/storage/models"
)

type NewGroup struct {
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Devices []string `json:"devices"`
}

func (n *NewGroup) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"name":    []string{"required"},
		"info":    []string{},
		"devices": []string{"uuid_array"},
	}
}

func (n *NewGroup) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"name":    []string{"Descriptive name for group. Required"},
		"info":    []string{"Additional information about group. Not required"},
		"devices": []string{"Devices to be added to group. Must be array of uuid4s."},
	}
}

// ToGroup Return dto as group
func (n *NewGroup) ToGroup() *models.Group {
	g := &models.Group{
		Name: n.Name,
		Info: n.Info,
	}
	if len(n.Devices) == 0 {
		return g
	}
	g.Devices = make([]models.Device, len(n.Devices))
	for i, v := range n.Devices {
		g.Devices[i].ID = v
	}
	return g
}

type Group struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Info    string   `json:"info"`
	Devices []string `json:"devices"`
}

func GroupToDto(g *models.Group) *Group {
	n := &Group{
		Id:   g.ID,
		Name: g.Name,
		Info: g.Info,
	}

	n.Devices = make([]string, len(g.Devices))
	for i, v := range g.Devices {
		n.Devices[i] = v.ID
	}
	return n
}
