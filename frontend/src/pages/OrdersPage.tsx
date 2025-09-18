import { useState, useCallback } from "react";
import CreateOrder from "../components/Orders/CreateOrder";
import { OrderList } from "../components/Orders/OrderList";
import './OrdersPage.scss';

export default function OrdersPage() {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleOrderCreated = useCallback(() => {
    setRefreshTrigger(prev => prev + 1);
  }, []);

  return (
    <div className="orders-page">
      <h1 className="orders-page__title">Управление заказами</h1>
      <div className="orders-page__content">
        <CreateOrder onOrderCreated={handleOrderCreated} />
        <OrderList key={`orders-${refreshTrigger}`} />
      </div>
    </div>
  );
}