package repo

import (
	"context"
	"github.com/DapperBlondie/users-cars-systems/src/models"
	zerolog "github.com/rs/zerolog/log"
	"time"
)

const (
	CarsTable = `CREATE TABLE cars 
( id integer NOT NULL PRIMARY KEY autoincrement , number_plate varchar(31) NOT NULL , color varchar(15) NOT NULL , vin varchar(31) NOT NULL , owner_id integer NOT NULL , CONSTRAINT vin_idx UNIQUE ( vin ) , CONSTRAINT num_idx UNIQUE ( number_plate ) , FOREIGN KEY ( owner_id ) REFERENCES users( id ) ON DELETE CASCADE ON UPDATE CASCADE )`

	UsersTable = `CREATE TABLE users 
( id integer NOT NULL PRIMARY KEY autoincrement , com_name varchar(63) NOT NULL , sex boolean NOT NULL , birthday time NOT NULL DEFAULT CURRENT_TIME , password char(31) NOT NULL )`

	AllCarsAndUsers = `SELECT s.id, s.com_name, s.sex, s.birthday, s.password, r.id, r.number_plate, r.color, r.vin, r.owner_id 
FROM users s INNER JOIN cars r ON ( r.owner_id = s.id )`

	GetUserCarsById = `SELECT r.id, r.number_plate, r.color, r.vin, r.owner_id FROM users s INNER JOIN cars r ON ( r.owner_id = s.id ) WHERE s.id=?`
)

type ApiOpsInterface interface {
	CreateTables() error
	AddUser(user *models.Users) error
	AddCar(car *models.Cars) error
	UpdateUser(user *models.Users) error
	UpdateCar(car *models.Cars) error
	DeleteUser(userID int) error
	GetUserByID(userID int) (*models.Users, error)
	GetAllUsers(skip int, limit int) ([]*models.Users, error)
}

// CreateTables use for creating our tables at the beginning of the program
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
