package repo

import (
	"context"
	zerolog "github.com/rs/zerolog/log"
	"time"
)

const (
	CarsTable = `CREATE TABLE cars 
( id integer NOT NULL PRIMARY KEY autoincrement , number_plate varchar(31) NOT NULL , color integer NOT NULL , vin varchar(31) NOT NULL , owner_id integer NOT NULL , CONSTRAINT vin_idx UNIQUE ( vin ) , CONSTRAINT num_idx UNIQUE ( number_plate ) , FOREIGN KEY ( owner_id ) REFERENCES users( id ) ON DELETE CASCADE ON UPDATE CASCADE )`
	UsersTable = `CREATE TABLE users 
( id integer NOT NULL PRIMARY KEY autoincrement , com_name varchar(100) NOT NULL , sex boolean NOT NULL , birthday time NOT NULL DEFAULT CURRENT_TIME , password char(32) NOT NULL ) `
)

type ApiOpsInterface interface {
	CreateTables() error
}

func (d *DBHolder) CreateTables() error {
	err := d.PingingDB()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = d.DB.ExecContext(ctx, CarsTable)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}

	_, err = d.DB.ExecContext(ctx, UsersTable)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}

	return nil
}
