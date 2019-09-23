package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(JsonMiddleware)

	router.HandleFunc("/graphql", BuildGraphQLAPI)

	http.Handle("/", router)

	return router
}
