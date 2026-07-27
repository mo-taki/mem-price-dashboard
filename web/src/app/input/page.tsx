"use client";

import { useState, useEffect } from "react";

type Product = {
  product: string;
  category: string;
  priceType: string;
};

export default function Input() {
  const [products, setProducts] = useState<Product[]>([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/products")
      .then((res) => res.json())
      .then((data) => setProducts(data));
  }, []);
  return (
    <main>
      <select>
        {products.map((p) => (
          <option key={p.product} value={p.product}>
            {p.product}
          </option>
        ))}
      </select>
    </main>
  );
}
