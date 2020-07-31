package tasker

import "time"

type User struct {
	Name     string
	Projects []*Project
}

type Project struct {
	Columns     []*Column
	Owner       *User
	Name        string
	Description string
}

type Column struct {
	Position int
	Project  *Project
	Tasks    []*Task
	Name     string
}

type Task struct {
	Name        string
	Description string
	Comments    []*Comment
	Column      *Column
}

type Comment struct {
	Author  *User
	Text    string
	Created time.Time
}
