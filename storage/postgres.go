package storage

import (
	"context"
	"database/sql"
	"time"
	"treasure/internal/rest/model"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db}
}

func (s *PostgresStorage) GetNames(searchName, searchType string) (*model.Persons, error) {
	return nil, nil
}

func (s *PostgresStorage) IsEmpty() bool {
	q := `select exists(select 1 from sdn limit 1)`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exists bool
	s.db.QueryRowContext(ctx, q).Scan(&exists)

	return exists
}

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
