package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DapperBlondie/users-cars-systems/src/models"
	zerolog "github.com/rs/zerolog/log"
	"time"
)

const (
	CarsTable = `CREATE TABLE IF NOT EXISTS cars  
( id integer NOT NULL PRIMARY KEY autoincrement , number_plate varchar(31) NOT NULL , color varchar(15) NOT NULL , vin varchar(31) NOT NULL , owner_id integer NOT NULL , CONSTRAINT vin_idx UNIQUE ( vin ) , CONSTRAINT num_idx UNIQUE ( number_plate ) , FOREIGN KEY ( owner_id ) REFERENCES users( id ) ON DELETE CASCADE ON UPDATE CASCADE )`

	UsersTable = `CREATE TABLE IF NOT EXISTS users
( id integer NOT NULL PRIMARY KEY autoincrement , com_name varchar(63) NOT NULL , sex boolean NOT NULL , birthday time NOT NULL DEFAULT CURRENT_TIME , password char(255) NOT NULL )`

	GetUserCarsById = `SELECT r.id, r.number_plate, r.color, r.vin FROM users s INNER JOIN cars r ON r.owner_id = s.id WHERE s.id=?`
)

type ApiOpsInterface interface {
	CreateTables() error
	AddUser(user *models.Users) error
	AddCar(car *models.Cars) error
	UpdateUser(user *models.Users) error
	UpdateCar(car *models.Cars) error
	DeleteUser(userID int) error
	GetUserByID(userID int) (*models.Users, error)
	GetAllUsers(limit, offset int) ([]*models.Users, error)
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
	_, err = d.DB.ExecContext(ctx, UsersTable)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}

	_, err = d.DB.ExecContext(ctx, CarsTable)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return err
	}

	return nil
}

// AddUser use for adding user into db
func (d *DBHolder) AddUser(user *models.Users) error {
	err := d.PingingDB()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	birthDay, err := time.Parse("2006-07-02", user.BirthDay)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	stmtQ := `INSERT INTO users (com_name, sex, birthday, password) VALUES (?, ?, ?, ?)`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	_, err = d.DB.ExecContext(ctx, stmtQ,
		user.CompleteName, user.Sex, birthDay, user.Password)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	return nil
}

// DeleteUser use for deleting a user with its own ID
func (d *DBHolder) DeleteUser(userID int) error {
	err := d.PingingDB()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	stmtQ := `DELETE FROM users WHERE id=? `
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	_, err = d.DB.ExecContext(ctx, stmtQ, userID)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	return nil
}

// AddCar use for adding car into the db
func (d *DBHolder) AddCar(car *models.Cars) error {
	err := d.PingingDB()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	query := `SELECT EXISTS(SELECT * FROM users WHERE id=?);`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result := d.DB.QueryRowContext(ctx, query, car.OwnerID)

	var rs int
	err = result.Scan(&rs)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}
	if rs == 0 {
		zerolog.Error().Msg(fmt.Sprintf("%d There is no car with this id=%d\n", rs, car.OwnerID))
		return errors.New(fmt.Sprintf("There is no car with this id=%d\n", car.OwnerID))
	}

	query = `INSERT INTO cars (number_plate,color,vin,owner_id) VALUES (?,?,?,?)`
	_, err = d.DB.ExecContext(ctx, query,
		car.NumberPlate, car.Color, car.VIN, car.OwnerID)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	return nil
}

// GetUserByID use for getting models.Users information with models.Cars
func (d *DBHolder) GetUserByID(userID int) (*models.Users, error) {
	err := d.PingingDB()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return nil, err
	}

	var user *models.Users = &models.Users{}
	query := `SELECT id,com_name,sex,birthday FROM users WHERE id=?`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	result := d.DB.QueryRowContext(ctx, query, userID)
	err = result.Scan(&user.ID,
		&user.CompleteName,
		&user.Sex,
		&user.BirthDay,
	)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return nil, err
	}

	results, err := d.DB.QueryContext(ctx, GetUserCarsById, userID)
	defer func(results *sql.Rows) {
		err = results.Close()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
	}(results)

	if err != nil {
		zerolog.Error().Msg(err.Error())
		return nil, err
	}
	if results == nil {
		return user, nil
	}

	var cars []*models.Cars = []*models.Cars{}
	car := &models.Cars{}
	for results.Next() {
		err = results.Scan(&car.ID,
			&car.NumberPlate,
			&car.Color,
			&car.VIN,
		)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return nil, err
		}

		cars = append(cars, car)
	}

	user.UsersCars = cars

	return user, nil
}

// GetAllUsers use for getting all users and associated cars
func (d *DBHolder) GetAllUsers(limit, offset int) ([]*models.Users, error) {
	err := d.PingingDB()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*25)
	defer cancel()

	var users []*models.Users
	query := `SELECT id FROM users LIMIT ? OFFSET ?`
	results, err := d.DB.QueryContext(ctx, query, limit, offset)
	defer func(results *sql.Rows) {
		err := results.Close()
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return
		}
	}(results)

	if err != nil {
		zerolog.Error().Msg(err.Error())
		return nil, err
	}
	if results == nil {
		return nil, errors.New("there is no any data available about users")
	}

	user := &models.Users{}
	for results.Next() {
		err = results.Scan(
			&user.ID,
		)
		if err != nil {
			zerolog.Error().Msg(err.Error())
			return nil, err
		}

		user, err = d.GetUserByID(user.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// UpdateUser use for update a user
func (d *DBHolder) UpdateUser(user *models.Users) error {
	err := d.PingingDB()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `UPDATE users SET com_name=?,sex=?,birthday=?,password=? WHERE id=?`
	_, err = d.DB.ExecContext(ctx, query,
		user.CompleteName,
		user.Sex,
		user.BirthDay,
		user.Password,
		user.ID)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return err
	}

	return nil
}

func (d *DBHolder) UpdateCar(car *models.Cars) error {
	return nil
}
