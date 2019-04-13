package dtos

import (
	"github.com/thedevsaddam/govalidator"
	"github.com/tryffel/fusio/storage/models"
)

// DeviceInfo Encapsulates only informational data about device
type DeviceInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func DeviceToInfo(d *models.Device) *DeviceInfo {
	return &DeviceInfo{
		Id:   d.ID,
		Name: d.Name,
		Type: d.DeviceType.ToString(),
	}
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Info string `json:"info"`
	Type string `json:"type"`
}

func FromDevice(device *models.Device) *Device {
	d := &Device{
		Id:   device.ID,
		Name: device.Name,
		Info: device.Info,
		Type: device.DeviceType.ToString(),
	}
	return d
}

type NewDevice struct {
	Name string `json:"name"`
	Info string `json:"info"`
	Type string `json:"type"`
}

func (n *NewDevice) ValidationMap() *govalidator.MapData {
	return &govalidator.MapData{
		"name": []string{"required", "min:1"},
		"info": []string{},
		"type": []string{"in:sensor,controller"},
	}
}

func (n *NewDevice) ValidationMessages() *govalidator.MapData {
	return &govalidator.MapData{
		"name": []string{"Descriptive name for device"},
		"info": []string{"Additiooanl information about the device"},
		"type": []string{"Either 'sensor' or 'controller'"},
	}
}
