package api

import (
	"encoding/json"
	"karna/core"
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

func LambdaAllHandler(w http.ResponseWriter, r *http.Request) {
	response := core.Lambda.BuildLambdaTree()

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func AGWAllHandler(w http.ResponseWriter, r *http.Request) {
	response := core.AGW.BuildAGWTree()
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func EC2AllHandler(w http.ResponseWriter, r *http.Request) {
	response := core.EC2.BuildEC2Tree()

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
