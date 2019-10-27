package api

import (
	"context"
	"flag"
	"github.com/karnalab/karna/core"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const (
	host = "127.0.0.1"
	port = "8000"
)

func startServer(router *mux.Router) {
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	server := &http.Server{
		Handler:      router,
		Addr:         host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			core.LogErrorMessage(err.Error())
		}
	}()

	core.LogSuccessMessage("Completed")
	core.LogSuccessMessage("API is listening @ " + server.Addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	server.Shutdown(ctx)

	core.LogSuccessMessage("API is shutting down")

	os.Exit(0)
}
