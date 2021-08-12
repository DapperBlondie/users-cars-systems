package routes

import (
	"github.com/DapperBlondie/users-cars-systems/src/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func ApiRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(handlers.ApiConf.EnableCORS)
	mux.Get("/status", handlers.ApiConf.CheckStatus)
	mux.Get("/delete-user", handlers.ApiConf.DeleteUserHandler)
	mux.Get("/get-user/{user_id}", handlers.ApiConf.GetUserHandler)
	mux.Get("/get-all-users", handlers.ApiConf.GetAllUsersHandler)

	mux.Post("/add-user", handlers.ApiConf.AddUserHandler)
	mux.Post("/add-car", handlers.ApiConf.AddCarHandler)

	return mux
}
