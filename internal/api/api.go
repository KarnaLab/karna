package api

//Run => Start Karna API.
func Run() {
	router := initRouter()
	startServer(router)
}
