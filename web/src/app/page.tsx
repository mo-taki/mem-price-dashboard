"use client";

import { useState, useEffect } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
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
  const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "";

  const loadMarket = () => {
    fetch(`${API_BASE}/api/market`)
      .then((res) => res.json())
      .then((data) => setMemPrices(data));
  };

  useEffect(() => {
    fetch(`${API_BASE}/api/prices?code=285A0`)
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
    fetch(`${API_BASE}/api/products`)
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

    fetch(`${API_BASE}/api/market`, {
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
    fetch(`${API_BASE}/api/market?${params.toString()}`, {
      method: "DELETE",
    }).then(() => loadMarket());
  };

  return (
    <main className="min-h-screen bg-white p-6 text-zinc-900">
      <h2 className="mt-4 mb-2 text-lg font-semibold">株価チャート</h2>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={mergePrices()}>
          <XAxis dataKey="date" />
          <YAxis yAxisId="stock" orientation="left" />
          <YAxis yAxisId="market" orientation="right" />
          <Tooltip />
          <Legend />
          <Line
            yAxisId="stock"
            dataKey="adjClose"
            name="キオクシア株価"
            stroke="#ff0000"
          />
          <Line
            yAxisId="market"
            dataKey="NAND 512Gb TLC Wafer"
            name="NAND 512Gb TLC Wafer"
            stroke="#82ca9d"
            connectNulls
          />
          <Line
            yAxisId="market"
            dataKey="DDR5 16Gb (2Gx8) 4800/5600"
            name="DDR5 16Gb (2Gx8) 4800/5600"
            stroke="#537afa"
            connectNulls
          />
          {/* <Line
            yAxisId="market"
            dataKey="NAND 256Gb TLC Wafer"
            name="NAND 256Gb TLC Wafer"
            stroke="#9c37fa"
            connectNulls
          /> */}
        </LineChart>
      </ResponsiveContainer>

      <hr />

      <h2 className="mt-8 mb-2 text-lg font-semibold">メモリ価格入力</h2>
      <form
        onSubmit={handleSubmit}
        className="flex flex-wrap items-end gap-3 my-5"
      >
        <select
          value={selectedProduct}
          onChange={(e) => setSelectedProduct(e.target.value)}
          className="rounded border border-zinc-300 px-3 py-2"
        >
          <option value="">-</option>
          {products.map((p) => (
            <option key={p.product} value={p.product}>
              {p.product}
            </option>
          ))}
        </select>
        <label className="flex flex-col text-sm">
          日付
          <input
            type="date"
            value={inputDate}
            onChange={(e) => setInputDate(e.target.value)}
            className="rounded border border-zinc-300 px-2 py-2"
          />
        </label>
        <label className="flex flex-col text-sm">
          価格
          <input
            type="number"
            value={inputPrice}
            onChange={(e) => setInputPrice(e.target.value)}
            className="rounded border border-zinc-300 px-2 py-2"
          />
        </label>
        <button
          type="submit"
          className="rounded bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700"
        >
          確定
        </button>
      </form>

      <hr />

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
                  onClick={() => handleDelete(m.date, m.product)}
                  className="text-red-600 hover:underline"
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
