package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/damianopetrungaro/golog"
	"github.com/damianopetrungaro/golog/opentelemetry"
	_ "github.com/lib/pq"
	"github.com/organization/order-service"
	"github.com/organization/order-service/cmd/internal/repo/cache"
	"github.com/organization/order-service/cmd/internal/repo/instrument"
	"github.com/organization/order-service/cmd/internal/repo/postgres"
	"github.com/organization/order-service/internal"
	"log"
	"os"
	"time"
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

	logger, flusher := newLogger()
	defer func() {
		flusher.Flush()
	}()

	db := newDB(ctx, logger)
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
		logger.With(golog.String("oder", fmt.Sprintf("%v", o))).Debug(ctx, "order was placed")
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
		logger.With(golog.String("oder", fmt.Sprintf("%v", o))).Debug(ctx, "order was shipped")
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
		logger.With(golog.String("oder", fmt.Sprintf("%v", o))).Debug(ctx, "order was delviered")
	default:
		logger.Error(ctx, "action not valid")
	}

	flusher.Flush()
	os.Exit(code)
}

func newLogger() (golog.Logger, golog.Flusher) {
	lvl, err := golog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatalf("could not parse log level: %s", err)
	}

	return opentelemetry.NewProductionLogger(lvl)
}

func newDB(ctx context.Context, logger golog.Logger) *sql.DB {
	const driver = "postgres"
	db, err := sql.Open(driver, os.Getenv("DB_URL"))
	if err != nil {
		logger.With(golog.Err(err)).Fatal(ctx, "could not connect to database")
	}

	return db
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
