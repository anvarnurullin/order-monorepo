import React, { useEffect, useState } from "react";
import type { Order } from "../types/order";
import { fetchOrders } from "../api/orders";
import "./OrderList.scss";

export const OrderList: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    fetchOrders().then(setOrders);
  }, []);

  return (
    <div className="order-list">
      <h2 className="order-list__title">Заказы</h2>
      <ul className="order-list__list">
        {orders.sort((a, b) => a.id - b.id ).map((o) => (
          <li key={o.id} className="order-list__item">
            Order #{o.id} - Product: {o.productId} - Qty: {o.quantity} - Status: {o.status}
          </li>
        ))}
      </ul>
    </div>
  );
};
