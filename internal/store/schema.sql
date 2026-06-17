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