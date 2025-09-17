import React, { useEffect, useState } from "react";
import type { Product } from "../types/product";
import { fetchProducts } from "../api/products";

export const ProductList: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);

  useEffect(() => {
    fetchProducts().then(setProducts);
  }, []);

  return (
    <div>
      <h2>Products</h2>
      <ul>
        {products.sort((a, b) => a.id - b.id ).map((p) => (
          <li key={p.id}>
            {p.id} - {p.name} (SKU: {p.sku}) - ${p.price} - Qty: {p.qtyAvailable}
          </li>
        ))}
      </ul>
    </div>
  );
};
