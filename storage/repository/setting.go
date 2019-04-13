package repository

type Setting interface {
	Get(key string) string
	Set(key string, value string)
}
