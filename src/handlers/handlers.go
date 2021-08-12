package handlers

import (
	"encoding/json"
	"errors"
	"github.com/DapperBlondie/users-cars-systems/src/models"
	"github.com/DapperBlondie/users-cars-systems/src/repo"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	zerolog "github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/scrypt"
	"net/http"
	"reflect"
	"strconv"
)

type ApiConfig struct {
	ScsManager *scs.SessionManager
	DHolder    *repo.DBHolder
}

var ApiConf *ApiConfig

func NewApiConf(scs *scs.SessionManager, dh *repo.DBHolder) {
	ApiConf = &ApiConfig{
		ScsManager: scs,
		DHolder:    dh,
	}
}

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

// CheckStatus just for showing the status of app
func (ac *ApiConfig) CheckStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		zerolog.Error().Msg(r.Method + " is not available")
		return
	}

	stat := &models.StatusIdentifier{
		Ok:      true,
		Message: "Everything is alright",
	}

	err := dResponseWriter(w, stat, http.StatusOK)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

// AddUserHandler use for adding users into the db
func (ac *ApiConfig) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		zerolog.Error().Msg(r.Method + " is not available")
		return
	}

	var user *models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPass)

	err = ac.DHolder.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stat := &models.StatusIdentifier{
		Ok:      true,
		Message: "User Added",
	}

	err = dResponseWriter(w, stat, http.StatusOK)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

// DeleteUserHandler use for deleting users from database
func (ac *ApiConfig) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		zerolog.Error().Msg(r.Method + " is not available")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is empty, fill it ", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "user_id is not an integer", http.StatusInternalServerError)
		return
	}

	err = ac.DHolder.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stat := &models.StatusIdentifier{
		Ok:      true,
		Message: "User Deleted",
	}

	err = dResponseWriter(w, stat, http.StatusOK)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

// AddCarHandler use for adding cars associated with user into db
func (ac *ApiConfig) AddCarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		zerolog.Error().Msg(r.Method + " is not available")
		return
	}

	var car *models.Cars
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ac.DHolder.AddCar(car)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stat := &models.StatusIdentifier{
		Ok:      true,
		Message: "Car Added",
	}

	err = dResponseWriter(w, stat, http.StatusOK)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}

// GetUserHandler use for get a user by its ID with associated cars
func (ac *ApiConfig) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		zerolog.Error().Msg(r.Method + " is not available")
		return
	}

	userID := chi.URLParamFromCtx(r.Context(), "user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := ac.DHolder.GetUserByID(id)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = dResponseWriter(w, user, http.StatusOK)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// GetAllUsersHandler get everything we have in db by limit & offset
func (ac *ApiConfig) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	lmt := r.URL.Query().Get("limit")
	off := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(lmt)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	offset, err := strconv.Atoi(off)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := ac.DHolder.GetAllUsers(limit, offset)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = dResponseWriter(w, users, http.StatusOK)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}
