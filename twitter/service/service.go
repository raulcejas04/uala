package service

import (
	database "twitter/models/database"
)

type DbService struct {
	db database.Database
}

// NewDbService creates a new instance of MyService
func NewDbService(db database.Database) *DbService {
	return &DbService{db: db}
}

// FetchData fetches data using the injected database
func (s *DbService) GetSeguidores(args ...interface{}) ([]map[string]interface{}, error) {
	return s.db.GetSeguidores(args...)
}
