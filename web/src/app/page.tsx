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

export default function Home() {
  const [prices, setPrices] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/prices?code=285A0")
      .then((res) => res.json())
      .then((data) => setPrices(data));
  }, []);

  return (
    <main>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={prices}>
          <XAxis dataKey="date" />
          <YAxis />
          <Tooltip />
          <Line dataKey="adjClose" />
        </LineChart>
      </ResponsiveContainer>
    </main>
  );
}
