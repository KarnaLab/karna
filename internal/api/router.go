package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(jsonMiddleware)
	router.Use(loggerMiddleware)

	router.HandleFunc("/graphql", buildGraphQLAPI)

	http.Handle("/", router)

	return router
}
