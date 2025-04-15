package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // blank import for sqlite driver
	"github.com/vinit-jpl/students-api-go/internal/config"

	// this is required for the sqlite driver to register itself with the database/sql package
	"github.com/vinit-jpl/students-api-go/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) { // create a new sqlite connection
	db, err := sql.Open("sqlite3", cfg.StoragePath) // returning a value and then error (common practise in go)
	if err != nil {
		return nil, err

	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil // returning a pointer to the Sqlite struct
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)") // passing "?" to prevent sql injection

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId() // this will return the last inserted id and error

	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1") // passing "?" to prevent sql injection

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age) // this will return the student with the given id

	if err != nil {
		if err == sql.ErrNoRows { // if no rows are found, return an empty student struct and nil error
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id)) // wrapping the error with additional context
		}
		return types.Student{}, fmt.Errorf("query error: %w", err) // wrapping the error with additional context)

	}

	return student, nil // returning the student struct and nil error
}
