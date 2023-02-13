package postgres

import (
	"context"
	"database/sql"
	"github.com/damianopetrungaro/golog"
	"github.com/google/uuid"
	"github.com/organization/order-service"
	"github.com/organization/order-service/cmd/internal/repo/postgres/internal"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	_ order.Repo = &Postgres{}
)

// Postgres represents a database layer for the order.Repo
type Postgres struct {
	db     *sql.DB
	logger golog.Logger
}

// New returns a database integration layer implementing order.Repo
func New(db *sql.DB, logger golog.Logger) *Postgres {
	return &Postgres{
		db:     db,
		logger: logger,
	}
}

// Find queries an order from the database
func (p *Postgres) Find(ctx context.Context, id order.ID) (*order.Order, error) {
	model, err := internal.Orders(
		qm.Where("id=?", id.String()),
	).One(ctx, p.db)
	if err != nil {
		p.logger.With(golog.Err(err)).Error(ctx, "order was not read from the database")
		return nil, order.ErrNotFound
	}

	return fromOrderModel(model), nil
}

// Add inserts an order to the database
func (p *Postgres) Add(ctx context.Context, o *order.Order) error {
	if err := toOrderModel(o).Upsert(ctx, p.db, true, []string{"id"}, boil.Infer(), boil.Infer()); err != nil {
		p.logger.With(golog.Err(err)).Error(ctx, "order was not inserted in the database")
		return order.ErrNotAdded
	}

	return nil
}

func fromOrderModel(model *internal.Order) *order.Order {
	return &order.Order{
		ID:          order.ID(uuid.MustParse(model.ID)),
		Number:      order.Number([]byte(model.Number)),
		Status:      order.Status(model.Status),
		PlacedBy:    order.UserID(uuid.MustParse(model.PlacedBy)),
		PlacedAt:    model.PlacedAt,
		ShippedAt:   model.ShippedAt.Time,
		DeliveredAt: model.DeliveredAt.Time,
	}
}

func toOrderModel(o *order.Order) *internal.Order {
	return &internal.Order{
		ID:          o.ID.String(),
		Number:      o.Number.String(),
		Status:      o.Status.String(),
		PlacedBy:    o.PlacedBy.String(),
		PlacedAt:    o.PlacedAt,
		ShippedAt:   null.TimeFrom(o.ShippedAt),
		DeliveredAt: null.TimeFrom(o.DeliveredAt),
	}
}
