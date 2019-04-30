package repository

// Cache acts as fast key-value store
type Cache interface {
	// Put value to cache. Value can be any primitive value
	// Set timeout to 0 to disable
	Put(key string, value interface{}, timeoutMs int) error

	// Get value
	Get(key string) (interface{}, error)

	// Delete key
	Delete(key string) error
}
