package repo

import (
	"context"
	"database/sql"
	zerolog "github.com/rs/zerolog/log"
)

type DBHolder struct {
	DB         *sql.DB
	Statements map[string]*sql.Stmt
}

var dbh *DBHolder

func NewDriver(dsn string) (*DBHolder, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return nil, err
	}

	dbh = &DBHolder{
		DB: db,
	}

	return dbh, nil
}

func (d *DBHolder) PingingDB() error {
	err := d.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (d *DBHolder) Dispose() error {
	err := d.DB.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d *DBHolder) CreateStatement(ctx context.Context, name string, query string) error {
	stmt, err := d.DB.PrepareContext(ctx, query)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	d.Statements[name] = stmt

	return nil
}
