package handlers

import (
	"encoding/json"
	"errors"
	"github.com/DapperBlondie/users-cars-systems/src/repo"
	"github.com/alexedwards/scs/v2"
	zerolog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"reflect"
)

type StatusIdentifier struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

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
			log.Println(err.Error())
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
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	}

	return errors.New("we could not be able to support data type that you passed")
}

func (ac *ApiConfig) CheckStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" is not available", http.StatusInternalServerError)
		zerolog.Error().Msg(r.Method + " is not available")
		return
	}

	stat := &StatusIdentifier{
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
