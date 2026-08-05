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

type Product = {
  product: string;
  category: string;
  priceType: string;
};

export default function Home() {
  const [stockPrices, setStockPrices] = useState<StockPrice[]>([]);
  const [memPrices, setMemPrices] = useState<MarketPrice[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [selectedProduct, setSelectedProduct] = useState("");
  const [inputDate, setInputDate] = useState("");
  const [inputPrice, setInputPrice] = useState("");

  const FROM = "2026-01-01";

  const loadMarket = () => {
    fetch("http://localhost:8080/api/market")
      .then((res) => res.json())
      .then((data) => setMemPrices(data));
  };

  useEffect(() => {
    fetch("http://localhost:8080/api/prices?code=285A0")
      .then((res) => res.json())
      .then((data) => setStockPrices(data));
  }, []);

  useEffect(() => {
    loadMarket();
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

    return [...byDate.values()]
      .filter((d) => d.date >= FROM)
      .sort((a, b) => a.date.localeCompare(b.date));
  };

  useEffect(() => {
    fetch("http://localhost:8080/api/products")
      .then((res) => res.json())
      .then((data) => setProducts(data));
  }, []);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const body = {
      product: selectedProduct,
      date: inputDate,
      price: Number(inputPrice),
    };

    fetch("http://localhost:8080/api/market", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }).then(() => loadMarket());
  };

  const handleDelete = (date: string, product: string) => {
    if (!window.confirm(`${date} / ${product} を削除しますか?`)) return;

    const params = new URLSearchParams({
      date: date,
      product: product,
    });
    fetch(`http://localhost:8080/api/market?${params.toString()}`, {
      method: "DELETE",
    }).then(() => loadMarket());
  };

  return (
    <main>
      <h2 className="mt-8 mb-2 text-lg font-semibold">株価チャート</h2>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={mergePrices()}>
          <XAxis dataKey="date" />
          <YAxis yAxisId="stock" orientation="left" />
          <YAxis yAxisId="market" orientation="right" />
          <Tooltip />
          <Line yAxisId="stock" dataKey="adjClose" stroke="#8884d8" />
          <Line
            yAxisId="market"
            dataKey="NAND 512Gb TLC Wafer"
            stroke="#82ca9d"
            connectNulls
          />
        </LineChart>
      </ResponsiveContainer>
      <h2 className="mt-8 mb-2 text-lg font-semibold">メモリ価格入力</h2>
      <form onSubmit={handleSubmit}>
        <select
          value={selectedProduct}
          onChange={(e) => setSelectedProduct(e.target.value)}
        >
          <option value="">-</option>
          {products.map((p) => (
            <option key={p.product} value={p.product}>
              {p.product}
            </option>
          ))}
        </select>
        <label>
          日付:{" "}
          <input
            type="date"
            onChange={(e) => setInputDate(e.target.value)}
            value={inputDate}
          />
        </label>
        <label>
          価格:{" "}
          <input
            type="number"
            onChange={(e) => setInputPrice(e.target.value)}
            value={inputPrice}
          />
        </label>
        <button type="submit">確定</button>
      </form>
      <h2 className="mt-8 mb-2 text-lg font-semibold">メモリ価格</h2>
      <table className="w-full border-collapse text-sm">
        <thead>
          <tr className="border-b border-zinc-300 bg-zinc-100 text-left text-zinc-700">
            <th className="px-3 py-2 font-medium">日付</th>
            <th className="px-3 py-2 font-medium">製品名</th>
            <th className="px-3 py-2 text-right font-medium">価格</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {memPrices.map((m) => (
            <tr
              key={`${m.product}-${m.date}`}
              className="border-b border-zinc-200 hover:bg-zinc-50"
            >
              <td className="px-3 py-2">{m.date}</td>
              <td className="px-3 py-2">{m.product}</td>
              <td className="px-3 py-2 text-right tabular-nums">{m.price}</td>
              <td className="px-3 py-2 text-right">
                <button
                  className="text-red-600 hover:underline"
                  onClick={() => handleDelete(m.date, m.product)}
                >
                  削除
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </main>
  );
}
