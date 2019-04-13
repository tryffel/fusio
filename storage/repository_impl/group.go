package repository_impl

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/models"
	"github.com/tryffel/fusio/storage/repository"
	"strings"
)

type GroupRepository struct {
	db *gorm.DB
}

func (g *GroupRepository) AddDevices(group *models.Group, ids []string) error {
	// Check that group owner owns devices also
	d := &[]models.Device{}
	count := 0

	res := g.db.Where("owner_id = ? AND id IN (?)", group.OwnerId, ids).Find(&d).Count(&count)
	if count != len(ids) {
		return &Err.Error{Code: Err.Einvalid, Err: errors.New("some devices are not found")}
	}
	if res.Error != nil {
		return getDatabaseError(res.Error)
	}

	devices := make([]models.Device, len(ids))
	for i, v := range ids {
		devices[i].ID = v
	}

	err := g.db.Model(&group).Association("Devices").Append(&devices).Error
	return getDatabaseError(err)
}

func (g *GroupRepository) UserHasAccess(userId uint, groupId []string) (bool, error) {
	group := &[]models.Group{}
	count := 0
	res := g.db.Where("owner_id = ? AND id IN (?)", userId, groupId).Find(group).Count(&count)
	if res.Error != nil {
		return false, getDatabaseError(res.Error)
	}

	if count == len(groupId) {
		return true, nil
	}
	return false, nil
}

func (g *GroupRepository) Create(group *models.Group) error {
	id := group.ID

	err := g.db.Create(&group).Error
	err = getDatabaseError(err)
	if err == nil {
		return nil
	}

	// Gorm updates group.ID regardless of whether operation succeedes or not
	// in case of existing group insert original id
	if e, ok := err.(*Err.Error); ok {
		if e.Message == "Already exists" {
			group.ID = id
		}
	}
	return err
}

func (g *GroupRepository) Remove(group *models.Group) error {
	err := g.db.Delete(&group).Error
	return getDatabaseError(err)
}

func (g *GroupRepository) Update(group *models.Group) error {
	err := g.db.Save(*group).Error
	return getDatabaseError(err)
}

func (g *GroupRepository) FindById(id string) (*models.Group, error) {
	group := &models.Group{}
	res := g.db.Where("id = ?", id).First(&group)
	return group, getDatabaseError(res.Error)
}

func (g *GroupRepository) FindByOwner(id int) (*[]models.Group, error) {
	group := &[]models.Group{}
	res := g.db.Where("owner_id = ?", id).Find(&group).Limit(100)
	return group, getDatabaseError(res.Error)
}

func (g *GroupRepository) FindByOwnerAndId(owner uint, id string) (*models.Group, error) {
	group := &models.Group{}
	res := g.db.Where("id = ? AND owner_id = ?", id, owner).First(&group)
	return group, getDatabaseError(res.Error)
}

func (g *GroupRepository) SearchGroups(opts *repository.SearchGroupsOpts) (*[]models.Group, error) {
	q := g.db
	if opts.OwnerId > 0 {
		q = q.Where("owner_id = ?", opts.OwnerId)
	}
	if opts.Name != "" {
		arg := strings.ToLower(opts.Name)
		arg = strings.Join([]string{"%", arg, "%"}, "")
		q = q.Where("lower_name LIKE ?", arg)
	}
	// TODO: implemenet showing deleted groups & filtering by devices

	groups := &[]models.Group{}
	res := q.Find(&groups)
	return groups, getDatabaseError(res.Error)
}

func (g *GroupRepository) GetDevices(owner uint, id string) (*[]string, error) {
	query :=
		`SELECT device_id
		FROM "groups"
		JOIN groups_devices ON group_id = groups.id
		WHERE groups.owner_id = ? 
		AND group_id = ?`

	//result := &[]idRes{}
	rows, err := g.db.Raw(query, owner, id).Rows()
	defer rows.Close()

	if err != nil {
		return &[]string{}, getDatabaseError(err)
	}
	var e error
	output := &[]string{}
	var res string
	for rows.Next() {
		res = ""
		err = rows.Scan(&res)
		if err != nil {
			e = Err.Wrap(&e, "Failed to scan device_ids to array")
		}
		*output = append(*output, res)
	}
	return output, getDatabaseError(e)
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	repo := &GroupRepository{
		db: db,
	}
	return repo
}
