package store

import (
	"database/sql"
	_ "embed"

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

type Price struct {
	Code      string  `json:"code"`
	Date      string  `json:"date"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	AdjClose  float64 `json:"adjClose"`
	AdjFactor float64 `json:"adjFactor"`
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
	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) ListPrices(code string) ([]Price, error) {
	q := `SELECT code, date, open, high, low, close, volume, adj_close, adj_factor FROM stock_prices WHERE code = ?;`
	rows, err := s.db.Query(q, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []Price

	for rows.Next() {
		var p Price
		if err := rows.Scan(&p.Code, &p.Date, &p.Open, &p.High, &p.Low, &p.Close, &p.Volume, &p.AdjClose, &p.AdjFactor); err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil

}
