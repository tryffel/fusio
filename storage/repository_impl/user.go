package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"strings"
)

type User struct {
	db *gorm.DB
}

func (u *User) Create(user *models.User, password string) error {
	err := user.SetPassword(password)
	if err != nil {
		return err
	}
	res := u.db.Create(&user)
	return getDatabaseError(res.Error)
}

func (u *User) Update(user *models.User) error {
	err := u.db.Save(&user).Error
	return getDatabaseError(err)
}

func (u *User) Delete(user *models.User) error {
	return u.db.Delete(&user).Error
}

func (u *User) FindById(id int) (*models.User, error) {
	user := &models.User{}
	res := u.db.Where("id = ? AND deleted_at IS NULL", id).First(&user)
	return user, res.Error
}

func (u *User) FindByName(name string) (*models.User, error) {
	user := &models.User{}
	res := u.db.Where("lower_name = ? AND deleted_at IS NULL", strings.ToLower(name)).First(&user)
	return user, res.Error
}

func (u *User) UpdatePassword(user *models.User, password string) error {
	err := user.SetPassword(password)
	if err != nil {
		return err
	}
	return u.Update(user)
}

func NewUserRepository(db *gorm.DB) repository.User {
	return &User{
		db: db,
	}
}
