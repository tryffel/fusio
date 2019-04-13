package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/util"
	"strings"
	"time"
)

// User type
type User struct {
	gorm.Model
	Name      string `gorm:"not null, unique"`
	LowerName string `gorm:"not null, unique"`
	Email     string `gorm:"not null, unique"`
	Password  string `gorm:"not null"`
	LastSeen  time.Time
	IsActive  bool `gorm:"not null"`
	IsAdmin   bool `gorm:"not null"`
	CanDelete bool `gorm:"not null"`
}

func (u *User) BeforeCreate() error {
	u.LowerName = strings.ToLower(u.Name)
	return nil
}

func (u *User) SetPassword(password string) error {
	hash, err := util.GetPasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

func GetUserByName(db *gorm.DB, name string) (*User, error) {
	u := &User{}
	db.Where("lower_name = ?", strings.ToLower(name)).First(&u)
	if u.ID < 1 {
		return u, errors.New("User not found")
	}
	return u, nil
}
