import type { Product } from "../types/product";

export async function fetchProducts(): Promise<Product[]> {
  const res = await fetch("/api/v1/products");
  const data = await res.json();
  return data.map((p: any) => ({
    id: p.id,
    name: p.name,
    sku: p.sku,
    price: p.price,
    qtyAvailable: p.qty_available,
    imageUrl: p.image_url,
    createdAt: p.created_at,
  }));
}
