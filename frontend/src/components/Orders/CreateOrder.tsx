import { useState } from "react";
import { createOrder as apiCreateOrder } from "../api/orders";
import "./CreateOrder.scss";

interface CreateOrderProps {
  onOrderCreated?: () => void;
}

export default function CreateOrder({ onOrderCreated }: CreateOrderProps) {
  const [productId, setProductId] = useState<number | null>(null);
  const [quantity, setQuantity] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const createOrder = async () => {
    if (!productId || !quantity) {
      alert("Пожалуйста, заполните все поля");
      return;
    }

    setIsLoading(true);
    try {
      await apiCreateOrder(productId, quantity);

      setProductId(null);
      setQuantity(null);

      onOrderCreated?.();
    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="create-order">
      <h2 className="create-order__title">Создать заказ</h2>
      <div className="create-order__form">
        <input
          type="number"
          value={productId || ""}
          onChange={(e) =>
            setProductId(e.target.value ? Number(e.target.value) : null)
          }
          placeholder="Product ID"
          className="create-order__input"
          disabled={isLoading}
        />
        <input
          type="number"
          value={quantity || ""}
          onChange={(e) =>
            setQuantity(e.target.value ? Number(e.target.value) : null)
          }
          placeholder="Quantity"
          className="create-order__input"
          disabled={isLoading}
        />
        <button
          onClick={createOrder}
          disabled={isLoading}
          className="create-order__button"
        >
          {isLoading ? "Создание..." : "Создать"}
        </button>
      </div>
    </div>
  );
}
