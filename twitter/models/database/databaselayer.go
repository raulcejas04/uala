package database

type Database interface {
	Connect() error
	GetSeguidores(args ...interface{}) ([]map[string]interface{}, error)
}
