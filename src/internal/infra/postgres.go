package infra

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	Conn *sqlx.DB
}

func NewPostgres(dataSourceName string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Postgres{Conn: db}, nil
}

func (h *Postgres) Db() *sqlx.DB {
	return h.Conn
}

func (h *Postgres) Close() error {
	return h.Conn.Close()
}
