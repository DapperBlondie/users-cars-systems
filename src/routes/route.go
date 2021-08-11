package routes

import (
	"github.com/DapperBlondie/users-cars-systems/src/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func ApiRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/status", handlers.ApiConf.CheckStatus)

	return mux
}
