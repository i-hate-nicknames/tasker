package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Run() {
	db, err := sql.Open("pgx", "postgresql://root:pass@localhost:8035/tasker")
	defer db.Close()
	if err != nil {
		log.Fatal("failed to connect to db", err.Error())
	}
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal("failed to connect to db: ping")
	// }

	err = Migrate(db)
	if err != nil {
		log.Fatal("failed to migrate", err.Error())
	}

	err = InsertThings(db)
	if err != nil {
		log.Fatal("failed to create projects: ", err.Error())
	}

	// things, err := GetThings(db)
	// if err != nil {
	// 	log.Fatal("failed to get things", err.Error())
	// }

	// for _, thing := range things {
	// 	log.Printf("%v\n", thing)
	// }

}

type Project struct {
	ID          int
	Columns     []*Column
	Name        string
	Description string
}

type Column struct {
	Position int
	Project  *Project `json:"-"`
	Name     string
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS Project (
		id serial PRIMARY KEY,
		name varchar(20) NOT NULL,
		description text NOT NULL,
        created timestamp
        );
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS Col (
		id serial PRIMARY KEY,
		project_id INT,
		position smallint NOT NULL,
		name varchar(20) NOT NULL,
		created timestamp,
		CONSTRAINT fk_project
			FOREIGN KEY(project_id)
			REFERENCES Project(id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
        );
	`)
	if err != nil {
		return err
	}

	return err
}

func InsertThings(db *sql.DB) error {
	stmtProj, err := db.Prepare("INSERT INTO Project (name, description, created) VALUES($1, $2, NOW()) RETURNING id")
	if err != nil {
		return err
	}
	stmtCol, err := db.Prepare("INSERT INTO Col (project_id, position, name, created) VALUES($1, $2, $3, NOW())")
	if err != nil {
		return err
	}
	defer stmtProj.Close()
	defer stmtCol.Close()
	var lastInsertedID int
	for i := 0; i < 10; i++ {
		err := stmtProj.QueryRow(fmt.Sprintf("name %d", i), "").Scan(&lastInsertedID)
		if err != nil {
			return err
		}
		for j := 0; j < 3; j++ {
			_, err := stmtCol.Exec(lastInsertedID, j, fmt.Sprintf("col name %d", i*10+j))
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// func GetThings(db *sql.DB) ([]*Thing, error) {
// 	rows, err := db.Query("SELECT * FROM test")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var (
// 		id      int
// 		name    string
// 		created time.Time
// 	)
// 	things := make([]*Thing, 0)
// 	for rows.Next() {
// 		err := rows.Scan(&id, &name, &created)
// 		if err != nil {
// 			return nil, err
// 		}
// 		thing := &Thing{ID: id, Name: name, Created: created}
// 		things = append(things, thing)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return things, nil
// }
