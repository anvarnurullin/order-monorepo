package store

import (
	"context"
	"fmt"
	"order-monorepo/services/order/internal/model"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore() (*Store, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@postgres:5433/app?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &Store{db: pool}, err
}


func (s *Store) CreateOrder (ctx context.Context, o model.Order) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx,
		`INSERT INTO orders (product_id, quantity, created_at) VALUES($1, $2, $3) RETURNING id`,
		o.ProductID, o.Quantity, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) GetOrders (ctx context.Context) ([]model.Order, error) {
	rows, err := s.db.Query(ctx, `SELECT id, product_id, quantity, created_at FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.Quantity, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
