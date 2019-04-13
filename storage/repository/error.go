package repository

type Errors interface {
	// Translate various errors e.g. sql errors into user friendly errors
	// Resources: gets inserted in error message if needed, e.g. 'device/alarm not found'
	GetUserFriendlyError(err error, resource string) error
}
