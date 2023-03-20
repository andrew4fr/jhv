package storage

import (
	"context"
	"database/sql"
	"time"
	"treasure/internal/rest/model"
)

type MariaDBStorage struct {
	db *sql.DB
}

func NewMariaDBStorage(db *sql.DB) *MariaDBStorage {
	return &MariaDBStorage{db}
}

func (s *MariaDBStorage) GetNames(searchName, searchType string) (*model.Persons, error) {
	return nil, nil
}

func (s *MariaDBStorage) SaveEntry(uid int, firstName, lastName string) error {
	q := `insert into sdn (uid, first_name, last_name) values ($1, $2, $3) 
		on duplicate key update first_name = values(first_name), last_name = values(last_name)`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, q, uid, firstName, lastName)

	if err != nil {
		return err
	}

	return nil
}
