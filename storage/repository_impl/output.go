package repository_impl

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"time"
)

type output struct {
	db *gorm.DB
}

func (o *output) Create(output *models.Output) error {
	err := o.db.Create(&output).Error
	return getDatabaseError(err)
}

func (o *output) Update(output *models.Output) error {
	err := o.db.Save(&output).Error
	return getDatabaseError(err)
}

func (o *output) Remove(output *models.Output) error {
	return &Err.Error{Code: Err.Econflict, Err: errors.New("removing outputs not implemented")}
}

func (o *output) FindbyId(id string) (*models.Output, error) {
	out := &models.Output{}
	res := o.db.Where("id = ?", id).Find(&out)
	return out, getDatabaseError(res.Error)
}

func (o *output) FindbyOwnerAndId(ownerId uint, id string) (*models.Output, error) {
	out := &models.Output{}
	res := o.db.Where("id = ? and owner_id = ?", id, ownerId).Find(&out)
	return out, getDatabaseError(res.Error)
}

func (o *output) FindByOwner(id uint) (*[]models.Output, error) {
	out := &[]models.Output{}
	res := o.db.Where("owner_id = ?", id).Find(&out)
	return out, getDatabaseError(res.Error)
}

func (o *output) FindByAlarm(alarmId string, opts repository.OutputOpts) (*[]models.Output, error) {
	conditions := map[string]interface{}{}
	conditions["alarm_id"] = alarmId

	if opts.OnlyEnabled {
		conditions["enabled"] = true
	}
	if opts.OnFire {
		conditions["on_fire"] = true
	}
	if opts.OnClear {
		conditions["on_clear"] = true
	}
	if opts.OnError {
		conditions["on_error"] = true
	}

	out := &[]models.Output{}
	res := o.db.Where(conditions).Preload("OutputChannel").Find(&out)
	return out, getDatabaseError(res.Error)
}

func (o *output) MarkPushed(output *models.Output, success bool, message string) error {
	output.LastPushed = time.Now()
	err := o.Update(output)
	if err != nil {
		return err
	}

	history := models.OutputHistory{
		OutputId: output.ID,
		Success:  success,
		Message:  message,
	}
	res := o.db.Create(&history)
	return getDatabaseError(res.Error)
}

func NewOutputRepository(db *gorm.DB) repository.Output {
	return &output{
		db: db,
	}
}
