"use client";

import { useState, useEffect } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

type StockPrice = {
  date: string;
  adjClose: number;
};

type MarketPrice = {
  date: string;
  product: string;
  price: number;
  category: string;
  priceType: string;
};

export default function Home() {
  const [stockPrices, setStockPrices] = useState<StockPrice[]>([]);
  const [memPrices, setMemPrices] = useState<MarketPrice[]>([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/prices?code=285A0")
      .then((res) => res.json())
      .then((data) => setStockPrices(data));
  }, []);

  useEffect(() => {
    fetch("http://localhost:8080/api/market")
      .then((res) => res.json())
      .then((data) => setMemPrices(data));
  }, []);

  const mergePrices = () => {
    const byDate = new Map();

    for (const s of stockPrices) {
      byDate.set(s.date, { date: s.date, adjClose: s.adjClose });
    }

    for (const m of memPrices) {
      const entry = byDate.get(m.date) ?? { date: m.date };
      entry[m.product] = m.price;
      byDate.set(m.date, entry);
    }

    return [...byDate.values()].sort((a, b) => a.date.localeCompare(b.date));
  };

  return (
    <main>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={mergePrices()}>
          <XAxis dataKey="date" />
          <YAxis yAxisId="stock" orientation="left" />
          <YAxis yAxisId="market" orientation="right" />
          <Tooltip />
          <Line yAxisId="stock" dataKey="adjClose" stroke="#8884d8" />
          <Line
            yAxisId="market"
            dataKey="DDR5 16Gb"
            stroke="#82ca9d"
            connectNulls
          />
        </LineChart>
      </ResponsiveContainer>
    </main>
  );
}
