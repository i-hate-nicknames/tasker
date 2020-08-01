package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Run() {
	db, err := sql.Open("pgx", "postgresql://root:pass@localhost:8035/tasker")
	defer db.Close()
	if err != nil {
		log.Fatal("failed to connect to db", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to connect to db: ping")
	}

	err = Migrate(db)
	if err != nil {
		log.Fatal("failed to migrate", err.Error())
	}

	err = InsertThings(db)
	if err != nil {
		log.Fatal("failed to insert things", err.Error())
	}

	things, err := GetThings(db)
	if err != nil {
		log.Fatal("failed to get things", err.Error())
	}

	for _, thing := range things {
		log.Printf("%v\n", thing)
	}

}

type Thing struct {
	ID      int
	Name    string
	Created time.Time
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Test (
                       id serial PRIMARY KEY,
                       name varchar(20) NOT NULL,
                       created timestamp
                       );
               `)
	return err
}

func InsertThings(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO test (name, created) VALUES($1, NOW())")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := 0; i < 10; i++ {
		_, err = stmt.Exec(fmt.Sprintf("thing %d", i))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetThings(db *sql.DB) ([]*Thing, error) {
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		id      int
		name    string
		created time.Time
	)
	things := make([]*Thing, 0)
	for rows.Next() {
		err := rows.Scan(&id, &name, &created)
		if err != nil {
			return nil, err
		}
		thing := &Thing{ID: id, Name: name, Created: created}
		things = append(things, thing)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return things, nil
}
