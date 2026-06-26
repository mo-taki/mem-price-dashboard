CREATE TABLE IF NOT EXISTS stock_prices (
    code        TEXT    NOT NULL,   -- "285A0"
    date        TEXT    NOT NULL,   -- "2026-03-17" (ISO8601)
    open        REAL,
    high        REAL,
    low         REAL,
    close       REAL,
    volume      INTEGER,
    adj_close   REAL,               -- 分割調整後の終値（チャートはこれを使う）
    adj_factor  REAL,
    PRIMARY KEY (code, date)
);

CREATE TABLE IF NOT EXISTS products (
    product     TEXT    PRIMARY KEY,    -- 製品名
    category    TEXT    NOT NULL,       -- DRAM / NAND
    price_type  TEXT    NOT NULL        -- spot / contract
);

CREATE TABLE IF NOT EXISTS market_prices (
    date    TEXT    NOT NULL,
    product TEXT    NOT NULL REFERENCES products(product),  -- 外部キー
    price   REAL,
    PRIMARY KEY (date, product)
);
