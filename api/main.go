package main

import (
	"github.com/mundanelizard/envi/pkg/logger"
	"os"
	"sync"
)

type server struct {
	logger *logger.Logger
	config struct {
		port int
	}
	wg sync.WaitGroup
}

func main() {
	srv := &server{
		logger: logger.New(os.Stdout, logger.Info),
	}

	err := srv.serve()
	if err != nil {
		srv.logger.Fatal(err)
	}
}
