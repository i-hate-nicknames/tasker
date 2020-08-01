package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/i-hate-nicknames/tasker/tasker"
)

// todo: the types  here are probably wrong and we should return some db layer types instead

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

type MemoryDb struct {
	projects []*tasker.Project
}

func MakeMemoryDb() *MemoryDb {
	user := &tasker.User{}
	p1 := tasker.MakeProject(user, "p1", "p1 descr")
	p2 := tasker.MakeProject(user, "p2", "p2 descr")
	p1.ID = 1
	p2.ID = 2
	ps := []*tasker.Project{p1, p2}
	return &MemoryDb{projects: ps}
}

func (db *MemoryDb) GetProjects() ([]*tasker.Project, error) {
	return db.projects, nil
}

var ErrNotFound = errors.New("not found")

func (db *MemoryDb) GetProject(id int) (*tasker.Project, error) {
	for _, p := range db.projects {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("get project: %w", ErrNotFound)
}
func (db *MemoryDb) SaveProject(p *tasker.Project) error {
	return nil
}
func (db *MemoryDb) DeleteProject(id int) error {
	return nil
}
