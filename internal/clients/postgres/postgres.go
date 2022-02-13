package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "err with Open DB")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "err with ping DB")
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) Begin() *sqlx.Tx {
	return p.db.MustBegin()
}
