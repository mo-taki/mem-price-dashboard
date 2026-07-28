"use client";

import { useState, useEffect } from "react";

type Product = {
  product: string;
  category: string;
  priceType: string;
};

export default function Input() {
  const [products, setProducts] = useState<Product[]>([]);
  const [selectedProduct, setSelectedProduct] = useState("");
  const [inputDate, setInputDate] = useState("");
  const [inputPrice, setInputPrice] = useState("");

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
    });
  };

  return (
    <main>
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
    </main>
  );
}
