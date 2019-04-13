package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

type Permission string

// Group Logical group of devices
type Group struct {
	ID        string `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	LowerName string `gorm:"not null"`
	Info      string
	Owner     User
	OwnerId   uint      `gorm:"not null"`
	Devices   []Device  `gorm:"many2many:groups_devices"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// BeforeCreate hook that gets called when creating new instance
func (g *Group) BeforeSave() (err error) {
	g.ID = util.NewUuid()
	g.LowerName = strings.ToLower(g.Name)
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return
}

// BeforeUpdate hook that updates timestamp
func (g *Group) BeforeUpdate() (err error) {
	g.UpdatedAt = time.Now()
	return
}

// AfterFind Fill Devices map as empty array
func (g *Group) AfterFind() (err error) {
	if g.Devices == nil {
		g.Devices = []Device{}
	}
	return nil
}

func (g *Group) AddDevice(db *gorm.DB, d *Device) error {
	res := db.Model(&g).Association("Devices").Append(d)
	logrus.Debug("Devices in group ", g.Name, ": ", len(g.Devices))
	return res.Error
}

func (g *Group) LoadDevices(db *gorm.DB) error {
	devices := &[]Device{}
	res := db.Model(&g).Association("Devices").Find(&devices)
	if res.Error != nil {
		return res.Error
	}
	g.Devices = *devices
	return res.Error
}

func (g *Group) Exists(db *gorm.DB) bool {
	group := &Group{}
	db.Where("id = ?", g.ID).First(&group)
	if group.ID != "" {
		return true
	}
	return false
}

// GetDevicesIds get groups devices ids
func (g *Group) GetDeviceIds() *[]string {
	s := make([]string, len(g.Devices))
	for i, v := range g.Devices {
		s[i] = v.ID
	}
	return &s
}
