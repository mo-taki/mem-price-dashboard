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

const upsertProductSQL = `INSERT INTO products (product, category, price_type)
		VALUES (?, ?, ?) ON CONFLICT (product)
		DO UPDATE SET category=excluded.category, price_type=excluded.price_type`

const upsertMarketPriceSQL = `INSERT INTO market_prices (date, product, price)
		VALUES (?, ?, ?) ON CONFLICT (date, product)
		DO UPDATE SET price=excluded.price`

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

type Product struct {
	Product   string `json:"product"`
	Category  string `json:"category"`
	PriceType string `json:"priceType"`
}

type MarketPrice struct {
	Date    string  `json:"date"`
	Product string  `json:"product"`
	Price   float64 `json:"price"`
}

type ProductMarketPrice struct {
	Date      string  `json:"date"`
	Product   string  `json:"product"`
	Price     float64 `json:"price"`
	Category  string  `json:"category"`
	PriceType string  `json:"priceType"`
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
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

func (s *Store) UpsertProduct(p Product) error {
	_, err := s.db.Exec(upsertProductSQL,
		p.Product, p.Category, p.PriceType)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ListProducts() ([]Product, error) {
	q := `SELECT product, category, price_type FROM products`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Product, &p.Category, &p.PriceType); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Store) UpsertMarketPrice(m MarketPrice) error {
	_, err := s.db.Exec(upsertMarketPriceSQL,
		m.Date, m.Product, m.Price)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ListMarketPrices(product string) ([]ProductMarketPrice, error) {
	q := `SELECT m.date, m.product, m.price, p.category, p.price_type
	FROM market_prices m
	JOIN products p ON m.product = p.product
	WHERE m.product = ?;`
	rows, err := s.db.Query(q, product)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productMarketPrices []ProductMarketPrice

	for rows.Next() {
		var pm ProductMarketPrice
		if err := rows.Scan(&pm.Date, &pm.Product, &pm.Price, &pm.Category, &pm.PriceType); err != nil {
			return nil, err
		}
		productMarketPrices = append(productMarketPrices, pm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productMarketPrices, nil
}
