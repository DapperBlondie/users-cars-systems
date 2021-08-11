package main

import (
	"github.com/DapperBlondie/users-cars-systems/src/handlers"
	"github.com/DapperBlondie/users-cars-systems/src/repo"
	"github.com/DapperBlondie/users-cars-systems/src/routes"
	"github.com/alexedwards/scs/v2"
	zerolog "github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	HOST   = "localhost"
	PORT   = ":9090"
	DBNAME = "app_db.db"
)

var session *scs.SessionManager

func main() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	dbh, err := repo.NewDriver(DBNAME)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}

	handlers.NewApiConf(session, dbh)

	srv := &http.Server{
		Addr:              HOST + PORT,
		Handler:           routes.ApiRoutes(),
		ReadTimeout:       time.Second * 11,
		ReadHeaderTimeout: time.Second * 6,
		WriteTimeout:      time.Second * 7,
		IdleTimeout:       time.Second * 6,
	}

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

	return
}
