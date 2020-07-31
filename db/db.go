package db

import (
	"database/sql"

	"github.com/i-hate-nicknames/tasker/tasker"
)

type Database interface {
	GetProjects() ([]*tasker.Project, error)
	GetProject(id int) (*tasker.Project, error)
	SaveProject(*tasker.Project) error
	DeleteProject(id int) error
}

type SqlDb struct {
	db sql.DB
}

func (db *SqlDb) GetProjects() ([]*tasker.Project, error) {
	return nil, nil
}
func (db *SqlDb) GetProject(id int) (*tasker.Project, error) {
	return nil, nil
}
func (db *SqlDb) SaveProject(*tasker.Project) error {
	return nil
}
func (db *SqlDb) DeleteProject(id int) error {
	return nil
}
