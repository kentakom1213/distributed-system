package storage

import (
	"context"
	"database/sql"
	_ "embed"
	
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

type Storage struct {
	db *sql.DB
}

func Open(ctx context.Context, path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	
	s := &Storage{
		db: db,
	}
	
	if err := s.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	
	return s, nil
}

func (s *Storage) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schemaSQL)
	return err
}

func (s *Storage) Close() error {
	return s.db.Close()
}
