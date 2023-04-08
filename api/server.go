package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (srv *server) serve() error {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.config.port),
		Handler:      srv.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     log.New(srv.logger, "HTTP", 0),
	}

	se := make(chan error)

	go srv.shutdown(s, se)

	srv.logger.Info(fmt.Sprintf("Server listening on %s", s.Addr))

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-se
	if err != nil {
		return err
	}

	srv.logger.Info("Server stopped", s.Addr)

	return nil
}

func (srv *server) shutdown(http *http.Server, e chan error) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit

	srv.logger.Info("Shutting down server", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := http.Shutdown(ctx)
	if err != nil {
		e <- err
	}

	srv.logger.Info("Completing background tasks", http.Addr)

	srv.wg.Wait()
	e <- nil
}
