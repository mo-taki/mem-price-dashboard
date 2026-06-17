package store

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/mo-taki/mem-price-dashboard/internal/jquants"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

const upsertStockPriceSQL = `INSERT INTO stock_prices (code, date, open, high, low, close, volume, adj_close, adj_factor)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(code, date)
		DO UPDATE SET open=excluded.open, high=excluded.high, low=excluded.low, close=excluded.close, volume=excluded.volume, adj_close=excluded.adj_close, adj_factor=excluded.adj_factor`

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schemaSQL)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) UpsertPrice(p jquants.StockPrice) error {
	_, err := s.db.Exec(upsertStockPriceSQL,
		p.Code, p.Date, p.O, p.H, p.L, p.C, p.Vo, p.AdjC, p.AdjFactor)
	if err != nil {
		return err
	}
	fmt.Println("Inserted/Updated:", p.Date, p.AdjC)
	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
