import React, { useEffect, useState } from "react";
import type { Order } from "../types/order";
import { fetchOrders } from "../api/orders";

export const OrderList: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    fetchOrders().then(setOrders);
  }, []);

  return (
    <div>
      <h2>Orders</h2>
      <ul>
        {orders.map((o) => (
          <li key={o.id}>
            Order #{o.id} - Product: {o.productId} - Qty: {o.quantity} - Status: {o.status}
          </li>
        ))}
      </ul>
    </div>
  );
};
