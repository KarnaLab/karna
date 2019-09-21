package api

func Start() {
	router := initRouter()
	startServer(router)
}

func Stop() {}
