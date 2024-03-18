package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/damianopetrungaro/golog"
	"github.com/damianopetrungaro/golog/opentelemetry"
	_ "github.com/lib/pq"
	"github.com/organization/order-service"
	"github.com/organization/order-service/cmd/internal/repo/cache"
	"github.com/organization/order-service/cmd/internal/repo/instrument"
	"github.com/organization/order-service/cmd/internal/repo/postgres"
	"github.com/organization/order-service/internal"
)

// This is a really "quick" implementation
// Feel free to open PR to make this more structured for a better real-world CLI example :)
func main() {
	var action, id, number, userID string
	flag.StringVar(&action, "action", "", "place, order, deliver")
	flag.StringVar(&number, "number", "", "order number to use when placing an order")
	flag.StringVar(&userID, "user_id", "", "user id to use when placing an order")
	flag.StringVar(&id, "id", "", "order id to use when delivering/shipping an order")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger, flusher, err := newLogger()
	if err != nil {
		log.Fatalf("getting new logger: %s", err)
	}
	defer func() {
		if err = flusher.Flush(); err != nil {
			logger.With(golog.Err(err)).Error(ctx, "could not flush logger")
		}
	}()

	db, err := newDB()
	if err != nil {
		logger.With(golog.Err(err)).Fatal(ctx, "could not connect to database")
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.With(golog.Err(err)).Error(ctx, "could not close database")
		}
	}()
	repo := newRepo(db, logger)

	svc := internal.NewService(repo, logger)

	var code int
	switch action {
	case "place":
		uID, err := order.ParseUserID(userID)
		if err != nil {
			logger.With(golog.Err(err)).Error(ctx, "user id was not valid")
			code = 1
			break
		}
		o, err := svc.Place(ctx, order.GenerateNumber(), uID)
		if err != nil {
			logger.With(golog.Err(err)).Error(ctx, "order was not placed")
			code = 2
			break
		}
		logger.With(golog.String("order", fmt.Sprintf("%v", o))).Debug(ctx, "order was placed")
	case "ship":
		_id, err := order.ParseID(id)
		if err != nil {
			logger.With(golog.Err(err)).Error(ctx, "id was not valid")
			code = 1
			break
		}
		o, err := svc.MarkAsShipped(ctx, _id)
		if err != nil {
			logger.With(golog.Err(err)).Error(ctx, "order was not shipped")
			code = 2
			break
		}
		logger.With(golog.String("order", fmt.Sprintf("%v", o))).Debug(ctx, "order was shipped")
	case "deliver":
		_id, err := order.ParseID(id)
		if err != nil {
			logger.With(golog.Err(err)).Error(ctx, "id was not valid")
			code = 1
			break
		}
		o, err := svc.MarkAsDelivered(ctx, _id)
		if err != nil {
			logger.With(golog.Err(err)).Error(ctx, "order was not delviered")
			code = 2
			break
		}
		logger.With(golog.String("order", fmt.Sprintf("%v", o))).Debug(ctx, "order was delviered")
	default:
		logger.Error(ctx, "action not valid")
	}

	flusher.Flush()
	os.Exit(code)
}

func newLogger() (gl golog.Logger, gf golog.Flusher, err error) {
	lvl, err := golog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return gl, gf, err
	}
	gl, gf = opentelemetry.NewProductionLogger(lvl)

	return gl, gf, nil
}

func newDB() (*sql.DB, error) {
	const driver = "postgres"
	db, err := sql.Open(driver, os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newRepo(db *sql.DB, logger golog.Logger) order.Repo {
	return instrument.New(
		cache.New(
			instrument.New(
				postgres.New(db, logger),
				"postgres",
			),
			cache.DefaultStore(),
			logger,
		),
		"cache",
	)
}
