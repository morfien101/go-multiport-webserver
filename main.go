package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	singalChan := make(chan os.Signal, 1)
	signal.Notify(singalChan, syscall.SIGINT, syscall.SIGTERM)

	portList := []int{8080, 8443}
	serverList := make([]*httpServer, 0)
	srvErrChan := make(chan error, 2)

	for _, port := range portList {
		logger := log.New(os.Stdout, fmt.Sprintf("%d - ", port), 0)
		srv := newHTTPServer(logger, port)
		srv.addHandlers(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello World! - port", srv.port)
			},
		)
		srv.server.Handler = srv.mux
		srv.setAddress("0.0.0.0", port)

		go func() {
			srvErrChan <- srv.server.ListenAndServe()
		}()
	}

	go func() {
		select {
		case err := <-srvErrChan:
			fmt.Printf("Error! Running web server shutting down. Error: %s", err)
			singalChan <- syscall.SIGTERM
		}
	}()

	signal := <-singalChan
	log.Println("Got signal", signal)
	// Shutdown here
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()
	errChan := make(chan error, len(serverList))
	for _, s := range serverList {
		errChan <- s.server.Shutdown(ctx)
	}

	numMsgs := 0
	exitCode := 0
	for {
		if numMsgs == len(serverList) {
			break
		}
		select {
		case err := <-errChan:
			numMsgs++
			if err != nil {
				exitCode = 1
				log.Printf("There was a problem shutting down the server. Error: %s\n", err)
			}
		}
	}
	// Terminate here
	os.Exit(exitCode)
}

type httpServer struct {
	server *http.Server
	logger *log.Logger
	mux    *http.ServeMux
	port   int
}

func newHTTPServer(logger *log.Logger, port int) *httpServer {
	return &httpServer{
		logger: logger,
		server: &http.Server{},
		mux:    http.NewServeMux(),
		port:   port,
	}
}

func (h *httpServer) setAddress(addr string, port int) {
	h.server.Addr = fmt.Sprintf("%s:%d", addr, port)
}

func (h *httpServer) addHandlers(route string, handlefunc http.HandlerFunc) {
	h.mux.HandleFunc(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Printf("Got reuest on %s.\n", h.server.Addr)
		handlefunc(w, r)
	}))
}
