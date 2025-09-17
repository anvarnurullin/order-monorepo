import type { Order } from "../types/order";

export async function fetchOrders(): Promise<Order[]> {
  const res = await fetch("/api/v1/orders");
  const data = await res.json();
  return data.map((o: any) => ({
    id: o.id,
    productId: o.product_id,
    quantity: o.quantity,
    status: o.status,
    createdAt: o.created_at,
  }));
}

export async function createOrder(productId: number, quantity: number): Promise<Order> {
  const res = await fetch("/api/v1/orders", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ product_id: productId, quantity }),
  });
  const data = await res.json();
  return {
    id: data.id,
    productId: data.product_id,
    quantity: data.quantity,
    status: data.status,
    createdAt: data.created_at,
  };
}
