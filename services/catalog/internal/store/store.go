package store

import (
	"context"
	"fmt"
	"order-monorepo/services/catalog/internal/model"
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
		dsn = "postgres://app:app@localhost:5433/app?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Store{db: pool}, err
}

func (s *Store) GetProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := s.db.Query(ctx, "SELECT id, name, sku, price, qty_available, created_at FROM products")
	if err != nil {
		return nil, err 
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Price, &p.QtyAvailable, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
