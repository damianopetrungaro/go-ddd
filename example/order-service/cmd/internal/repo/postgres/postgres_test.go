package postgres

import (
	"context"
	"database/sql"
	gologTest "github.com/damianopetrungaro/golog/test"
	"github.com/google/uuid"
	"github.com/organization/order-service"
	"github.com/organization/order-service/cmd/internal/repo/postgres/internal"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"testing"
	"time"
)

func TestPostgres_Add(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(func() {
		cancel()
	})

	o := getRandomOrder(t)
	db := getDB(t)
	repo := getPostgres(t, db)

	if err := repo.Add(ctx, o); err != nil {
		t.Fatalf("could not add order: %v", o)
	}

	found := getOrderByIDHelper(t, repo.db, o.ID.String())
	matchesOrder(t, o, found)
}

func TestOrderRepo_GetOrder(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(func() {
		cancel()
	})

	o := getRandomOrder(t)
	db := getDB(t)
	repo := getPostgres(t, db)

	addOrderHelper(t, repo.db, o)
	found, err := repo.Find(ctx, o.ID)
	if err != nil {
		t.Fatalf("could not get order: %v", o)
	}

	matchesOrder(t, o, found)
}

func getPostgres(t *testing.T, db *sql.DB) *Postgres {
	t.Helper()

	return New(db, gologTest.NewNullLogger())
}

func addOrderHelper(t *testing.T, db *sql.DB, b *order.Order) {
	if err := toOrderModel(b).Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Fatalf("could not execute insert query: %s", err)
	}
}

func getOrderByIDHelper(t *testing.T, db *sql.DB, id string) *order.Order {
	t.Helper()
	model, err := internal.Orders(qm.Where("id=?", id)).One(context.Background(), db)
	if err != nil {
		t.Fatalf("could not query order")
	}

	return fromOrderModel(model)
}

func getRandomOrder(t *testing.T) *order.Order {
	t.Helper()

	return &order.Order{
		ID:          order.ID(uuid.New()),
		Number:      order.GenerateNumber(),
		Status:      order.Shipped,
		PlacedBy:    order.UserID(uuid.New()),
		PlacedAt:    time.Now(),
		ShippedAt:   time.Now(),
		DeliveredAt: time.Time{},
	}
}

func matchesOrder(t *testing.T, a, b *order.Order) {
	defer func() {
		if t.Failed() {
			t.Error("could not match orders")
			t.Errorf("a: %#v", a)
			t.Errorf("b: %#v", b)
		}
	}()

	const timeLayout = "2006-01-02T15:04:05"

	if a.ID != b.ID {
		t.Fail()
	}
	if a.Number != b.Number {
		t.Fail()
	}
	if a.Status != b.Status {
		t.Fail()
	}
	if a.PlacedBy != b.PlacedBy {
		t.Fail()
	}
	if a.PlacedAt.Compare(a.PlacedAt) != 0 {
		t.Fail()
	}
	if a.ShippedAt.Compare(a.ShippedAt) != 0 {
		t.Fail()
	}
	if a.DeliveredAt.Compare(a.DeliveredAt) != 0 {
		t.Fail()
	}
}
