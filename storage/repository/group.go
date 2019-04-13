package repository

import "github.com/tryffel/fusio/storage/models"

type SearchGroupsOpts struct {
	OwnerId     uint
	Name        string
	DeviceIds   []string
	ShowDeleted bool
}

type Group interface {
	// CRUD
	Create(group *models.Group) error
	Remove(group *models.Group) error
	Update(group *models.Group) error
	FindById(id string) (*models.Group, error)
	FindByOwner(id int) (*[]models.Group, error)
	FindByOwnerAndId(owner uint, id string) (*models.Group, error)

	// Check if user has access to given groups
	UserHasAccess(userId uint, groupId []string) (bool, error)

	// AddDevices adds given devices to group. This doesn't validate devices
	AddDevices(group *models.Group, ids []string) error
	// Search groups by given options.
	SearchGroups(opts *SearchGroupsOpts) (*[]models.Group, error)

	//GetDevices returns devices assigned to group
	GetDevices(owner uint, group string) (*[]string, error)
}
