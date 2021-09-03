package web

import (
	"currency-conversion/config"
	"currency-conversion/web/handlers"
	"currency-conversion/web/middlewares"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start() {

	http.Handle("/convert", http.HandlerFunc(handlers.HandlerConvert))
	http.Handle("/create", http.HandlerFunc(middlewares.Authorize(handlers.HandlerCreateProject)))

	cfg := config.Get()

	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan error)

	//Close server on signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//Wait for signal and close server
	go func() {
		<-sigs
		server.Close()
	}()

	//wait for server to be closed
	go func() {
		err := server.ListenAndServe()
		done <- err
	}()

	//Exit when server is closed
	err := <-done

	if err != nil {
		log.Fatal(err)
	}
}
