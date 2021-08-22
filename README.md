# users-cars-systems
 A system for managing users with its cars that have limited functionalities
 
 ## Main Function
 I define a functionality for exiting the program peacfully.
 
 ```go
 
 sigC := make(chan os.Signal, 1)
	signal.Notify(sigC)

	go func() {
		zerolog.Log().Msg("HTTP1.x server is listening on " + HOST + PORT)
		if err := srv.ListenAndServe(); err != nil {
			zerolog.Fatal().Msg(err.Error())
			return
		}
	}()

	<-sigC
 
 ```
 
 ***
 
 ## Data Models
 I create three models for send response to the client.
 - First ``` StatusIdentifier ``` for writing a status of operation.
 - Second ``` Users ``` for handling payloads and datas that associated with Users.
 - Third ``` Cars ``` for handling payloads and datas that associated with Cars.
 
 ```go
 package models

type StatusIdentifier struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// Users holding users data that stored in DB in a structures
type Users struct {
	ID           int     `json:"id,omitempty"`
	CompleteName string  `json:"complete_name"`
	Sex          bool    `json:"sex"`
	BirthDay     string  `json:"birth_day"`
	Password     string  `json:"password"`
	UsersCars    []*Cars `json:"users_cars,omitempty"`
}

// Cars holding cars data
type Cars struct {
	ID          int    `json:"id,omitempty"`
	NumberPlate string `json:"number_plate"`
	Color       string `json:"color"`
	VIN         string `json:"vin"`
	OwnerID     int    `json:"owner_id"`
}

 ```
 
 ***
 
 ## SQL Queries
 I store some repeated queries in consts variables.You can find them in ``` repository.go ``` .
 
 - First a SQL query that runs at the beginning of the program for creating the cars table.
 
 ```sql
 CREATE TABLE IF NOT EXISTS cars  
( id integer NOT NULL PRIMARY KEY autoincrement , number_plate varchar(31) NOT NULL , color varchar(15) NOT NULL , vin varchar(31) NOT NULL , owner_id integer NOT NULL , CONSTRAINT vin_idx UNIQUE ( vin ) , CONSTRAINT num_idx UNIQUE ( number_plate ) , FOREIGN KEY ( owner_id ) REFERENCES users( id ) ON DELETE CASCADE ON UPDATE CASCADE )
 ```

- Second a SQL query that runs at the beginning of the program for creating the users table.

```sql
CREATE TABLE IF NOT EXISTS users
( id integer NOT NULL PRIMARY KEY autoincrement , com_name varchar(63) NOT NULL , sex boolean NOT NULL , birthday time NOT NULL DEFAULT CURRENT_TIME , password char(255) NOT NULL )
```

- Third a SQL query that use for get a user by its own ID from DataBase with associated cars.

```sql
SELECT r.id, r.number_plate, r.color, r.vin FROM users s INNER JOIN cars r ON r.owner_id = s.id WHERE s.id=?

```

***

## Sqlite3 Driver
I write a simple driver for managing our connection with DB.

- I define a structure for holding the database connection and other probable utilities.

```go
type DBHolder struct {
	DB         *sql.DB
	Statements map[string]*sql.Stmt
}

```

- I write a ``` Dispose ``` function for release the associated data and memory at the end of program.

```go
func (d *DBHolder) Dispose() error {
	err := d.DB.Close()
	if err != nil {
		return err
	}

	return nil
}

```

- I define an ``` interface ``` that can be used for better ordering of DB; Because, it has some basic functionalities that can implements by every RDBMS.

```go

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

```

***

## Hashing and Encrypting Password
I am so sorry for this part of task becasue of **PRIVACY POLICY**; We and other people except than user **SHOULD NOT BE ABLE TO SEE THE USER PASSWORD** and decrypt it.<br/>
- For Encrypting, I use this ``` golang.org/x/crypto/bcrypt ``` package; But, for high privacy you can use ``` golang.org/x/crypto/srypt ``` and for cost I use 12.

```go 

hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPass)
	
```

- For Decrypting, I can do it because of **PRIVACY POLICY**. For changing it we can define its handler like twitter or linkedin. 
- For comparing the password I use this functionality.

```go 

err := bcrypt.CompareHashAndPassword(hashedPass, normalPass)
	if err != nil {
		return 
	}

```

## ResponseWriter
I wrote a reponse writer for [web-auth-methods](https://gist.github.com/DapperBlondie/872ffeea7da05a600d93a78f00ebe2e4) project.
But for this project I applied some modification. I add support for writing list of objects; Actually last ``` else-if ``` .

```go

// dResponseWriter use for writing response to the user
func dResponseWriter(w http.ResponseWriter, data interface{}, HStat int) error {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.String {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/text")

		_, err := w.Write([]byte(data.(string)))
		return err
	} else if reflect.PtrTo(dataType).Kind() == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			zerolog.Error().Msg(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Struct == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			zerolog.Error().Msg(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Slice == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			zerolog.Error().Msg(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	}

	return errors.New("we could not be able to support data type that you passed")
}

```

***

## API Operations
At first I show URLs that I use for API it can define the schema of URLs manifestly.

```url
http://localhost:9090/add-user

http://localhost:9090/add-car

http://localhost:9090/get-user/{user_id}

http://localhost:9090/get-all-users?limit=<integer_numbet>&offset=<integer_numbet>

http://localhost:9090/update-user

http://localhost:9090/update-car

http://localhost:9090/delete-user?user_id=1
```

### GetUserHandler
I get user_id from url then I returned back the user.

- Also I use ``` context.WithTimeout ``` for control the time of executions.


```go

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

```

### GetAllUsersHandler
I get two parameter ``` limit ``` , ``` offset ``` for specify the amount of data we want to show.

- Also I use ``` context.WithTimeout ``` for control the time of executions.

```go 

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

```
