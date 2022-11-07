package infra

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MariaDB struct {
	Conn *sqlx.DB
}

func NewMariaDB(dataSourceName string) (*MariaDB, error) {
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &MariaDB{Conn: db}, nil
}

func (h *MariaDB) Db() *sqlx.DB {
	return h.Conn
}

func (h *MariaDB) Close() error {
	return h.Conn.Close()
}
