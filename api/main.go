package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/mundanelizard/envi/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type server struct {
	logger *logger.Logger
	config struct {
		port int
	}
	wg  sync.WaitGroup
	db  *mongo.Database
	ctx context.Context
}

func main() {
	srv := &server{
		logger: logger.New(os.Stdout, logger.Info),
		ctx:    context.Background(),
	}

	srv.config.port = 6000

	db, err := loadDb()
	if err != nil {
		srv.logger.Fatal(err)
	}

	srv.db = db

	err = srv.serve()
	if err != nil {
		srv.logger.Fatal(err)
	}
}

func loadDb() (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("envi")

	return db, nil
}
