package storage

import (
	"context"
	"database/sql"
	"strings"
	"time"
	"treasure/internal/rest/model"
)

type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage makes new storage instance
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db}
}

// StrongGetNames finds person by exactly match
func (s *PostgresStorage) StrongGetNames(searchName string) (*model.Persons, error) {
	q := `select uid, first_name, last_name from sdn
		where lower(first_name) = lower($1) or lower(last_name) = lower($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, q, searchName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result model.Persons
	for rows.Next() {
		var person model.Person
		err := rows.Scan(&person.UID, &person.FirstName, &person.LastName)
		if err != nil {
			return nil, err
		}
		result = append(result, &person)
	}

	return &result, nil
}

// StrongGetNames finds person by non exactly match
func (s *PostgresStorage) WeakGetNames(searchName string) (*model.Persons, error) {
	q := `select uid, first_name, last_name from sdn
		where lower(first_name) like lower($1) or lower(last_name) like lower($1)`

	stmt, err := s.db.Prepare(q)
	if err != nil {
		return nil, err
	}

	var result model.Persons
	names := strings.Split(searchName, " ")
	for _, name := range names {
		res, err := s.weakGetName(name, stmt)
		if err != nil {
			return nil, err
		}

		result = append(result, *res...)
	}

	uniqueResult := unique(&result)

	return uniqueResult, nil
}

func (s *PostgresStorage) weakGetName(name string, stmt *sql.Stmt) (*model.Persons, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := stmt.QueryContext(ctx, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result model.Persons
	for rows.Next() {
		var person model.Person
		err := rows.Scan(&person.UID, &person.FirstName, &person.LastName)
		if err != nil {
			return nil, err
		}
		result = append(result, &person)
	}

	return &result, nil

}

// IsEmpty checks for data presenting
func (s *PostgresStorage) IsEmpty() bool {
	q := `select exists(select 1 from sdn limit 1)`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exists bool
	s.db.QueryRowContext(ctx, q).Scan(&exists)

	return exists
}

// SaveEntry saves person data
func (s *PostgresStorage) SaveEntry(uid int, firstName, lastName string) error {
	q := `insert into sdn (uid, first_name, last_name) values ($1, $2, $3) 
		on conflict (uid) do update set 
		first_name = excluded.first_name, last_name = excluded.last_name`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, q, uid, firstName, lastName)

	if err != nil {
		return err
	}

	return nil
}

func unique(persons *model.Persons) *model.Persons {
	var result model.Persons
	uids := make(map[int64]bool)

	for _, p := range *persons {
		if _, ok := uids[p.UID]; !ok {
			uids[p.UID] = true
			result = append(result, p)
		}
	}

	return &result
}
