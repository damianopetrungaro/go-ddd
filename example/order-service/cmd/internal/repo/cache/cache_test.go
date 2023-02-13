package cache

import (
	"context"
	"errors"
	"github.com/damianopetrungaro/go-cache"
	gologTest "github.com/damianopetrungaro/golog/test"
	"github.com/golang/mock/gomock"
	"github.com/organization/order-service"
	"testing"
)

func TestCache_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		o := &order.Order{ID: order.NewID()}

		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		repo.EXPECT().Add(ctx, o).Times(1).Return(nil)

		cachedRepo := New(repo, DefaultStore(), logger)

		if err := cachedRepo.Add(ctx, o); err != nil {
			t.Fatalf("could not add: %s", err)
		}

		found, err := cachedRepo.store.Get(ctx, o.ID)
		if err != nil {
			t.Fatalf("could not find order in the store: %s", err)
		}

		if found != o {
			t.Error("could not match orders")
			t.Errorf("got: %v", found)
			t.Errorf("want: %v", o)
		}
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		o := &order.Order{ID: order.NewID()}

		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		repo.EXPECT().Add(ctx, o).Times(1).Return(order.ErrNotAdded)

		cachedRepo := New(repo, DefaultStore(), logger)

		if err := cachedRepo.Add(ctx, o); !errors.Is(err, order.ErrNotAdded) {
			t.Fatalf("could not match error: %s", err)
		}

		_, err := cachedRepo.store.Get(ctx, o.ID)
		if !errors.Is(err, cache.ErrNotFound) {
			t.Fatalf("could find order in the store: %s", err)
		}
	})
}

func TestCache_Find(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		o := &order.Order{ID: order.NewID()}

		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		repo.EXPECT().Find(ctx, o.ID).Times(1).Return(o, nil)

		cachedRepo := New(repo, DefaultStore(), logger)

		found, err := cachedRepo.Find(ctx, o.ID)
		if err != nil {
			t.Fatalf("could not find order: %s", err)
		}

		if found != o {
			t.Error("could not match orders")
			t.Errorf("got: %v", found)
			t.Errorf("want: %v", o)
		}

		found, err = cachedRepo.Find(ctx, o.ID)
		if err != nil {
			t.Fatalf("could not find order: %s", err)
		}

		if found != o {
			t.Error("could not match orders")
			t.Errorf("got: %v", found)
			t.Errorf("want: %v", o)
		}
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		ctx := context.Background()
		o := &order.Order{ID: order.NewID()}

		repo := order.NewMockRepo(ctrl)
		logger := gologTest.NewNullLogger()

		repo.EXPECT().Find(ctx, o.ID).Times(2).Return(nil, order.ErrNotFound)

		cachedRepo := New(repo, DefaultStore(), logger)

		_, err := cachedRepo.Find(ctx, o.ID)
		if !errors.Is(err, order.ErrNotFound) {
			t.Fatalf("could find order: %s", err)
		}

		_, err = cachedRepo.Find(ctx, o.ID)
		if !errors.Is(err, order.ErrNotFound) {
			t.Fatalf("could find order: %s", err)
		}
	})
}
