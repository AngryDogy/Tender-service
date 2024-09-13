package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"tenderservice/logger"
)

type Repository interface {
	Connect() error
	Initialize(schemePath string) error
	GetConnection() *sql.DB
	Close() error
}

type postgresRepository struct {
	source string
	conn   *sql.DB
}

func NewPostgresRepository(source string) Repository {
	return &postgresRepository{
		source: source,
	}
}

func (r *postgresRepository) Connect() error {
	conn, err := sql.Open("postgres", r.source)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}
	r.conn = conn

	return nil
}

func (r *postgresRepository) Initialize(schemePath string) error {
	file, err := os.ReadFile(schemePath)
	if err != nil {
		return err
	}

	queries := strings.Split(string(file), ";")
	for _, query := range queries {
		_, err := r.conn.Exec(query)
		if err != nil {
			logger.WarnLogger.Println("failed to execute query", err)
		}
	}

	return nil
}

func (r *postgresRepository) GetConnection() *sql.DB {
	return r.conn
}

func (r *postgresRepository) Close() error {
	return r.conn.Close()
}
