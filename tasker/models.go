package tasker

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID       int
	Name     string
	Projects []*Project
}

type Project struct {
	ID          int
	Columns     []*Column
	Owner       *User
	Name        string
	Description string
}

type Column struct {
	Position int
	Project  *Project `json:"-"`
	Tasks    []*Task
	Name     string
}

type Task struct {
	Author      *User
	Name        string
	Description string
	Comments    []*Comment
	Column      *Column
	Position    int
	Created     time.Time
}

type Comment struct {
	Author  *User
	Title   string
	Text    string
	Created time.Time
}

// MakeProject creates a project with specified owner and information
// and a single default column
func MakeProject(owner *User, name, description string) *Project {
	tasks := make([]*Task, 0)
	cols := []*Column{{Position: 0, Tasks: tasks}}
	p := &Project{Owner: owner, Columns: cols, Name: name, Description: description}
	cols[0].Project = p
	return p
}

// --- Columns

var ErrLastColumn = errors.New("trying to remove last column in a project")

// DeleteColumn from the given project at given position
// Do nothing if there is no column at the given position
// return ErrLastColumn error if this is the last column in the project
func DeleteColumn(p *Project, pos int) error {
	if len(p.Columns) == 1 {
		return ErrLastColumn
	}
	return nil
}

var ErrDifferentProjects = errors.New("entities belong to different projects")

// Swap making position of col to be position of target
// and vice versa
// return ErrDifferentProjects when columns belong to different projects
func (col *Column) Swap(target *Column) error {
	if col.Project != target.Project {
		return fmt.Errorf("swap columns: %w", ErrDifferentProjects)
	}
	col.Position, target.Position = target.Position, col.Position
	return nil
}

// --- Tasks

func MakeTask(author *User, column *Column, name, description string) *Task {
	now := time.Now()
	return &Task{Name: name, Description: description, Author: author, Created: now}
}

func (task *Task) MoveToColumn(target *Column) error {
	if task.Column.Project != target.Project {
		return fmt.Errorf("move to column: %w", ErrDifferentProjects)
	}
	if task.Column == target {
		return nil
	}
	newPos := 0
	if len(target.Tasks) > 0 {
		lastPos := target.Tasks[len(target.Tasks)-1].Position
		newPos = lastPos + 1
	}
	// todo: remove task from its original column
	task.Column = target
	task.Position = newPos
	target.Tasks = append(target.Tasks, task)
	return nil
}

var ErrDifferentColumns = errors.New("tasks belong to different columns")

func (task *Task) Swap(target *Task) error {
	if task.Column != target.Column {
		return fmt.Errorf("swap: %w", ErrDifferentColumns)
	}
	task.Position, target.Position = target.Position, task.Position
	return nil
}

// --- Comments

func MakeComment(author *User, task *Task, title, text string) *Comment {
	now := time.Now()
	return &Comment{Title: title, Text: text, Author: author, Created: now}
}
