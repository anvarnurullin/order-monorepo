import { useState } from "react";
import { createOrder as apiCreateOrder } from "../api/orders";

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

      // Clear form
      setProductId(null);
      setQuantity(null);

      // Trigger refresh of parent components
      onOrderCreated?.();

    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="p-4 border rounded-md shadow">
      <h2 className="text-xl font-bold mb-2">Создать заказ</h2>
      <div className="flex gap-2 mb-2">
        <input
          type="number"
          value={productId || ""}
          onChange={(e) =>
            setProductId(e.target.value ? Number(e.target.value) : null)
          }
          placeholder="Product ID"
          className="border p-1"
          disabled={isLoading}
        />
        <input
          type="number"
          value={quantity || ""}
          onChange={(e) =>
            setQuantity(e.target.value ? Number(e.target.value) : null)
          }
          placeholder="Quantity"
          className="border p-1"
          disabled={isLoading}
        />
        <button
          onClick={createOrder}
          disabled={isLoading}
          className="bg-blue-500 text-white px-3 py-1 rounded disabled:bg-gray-400"
        >
          {isLoading ? "Создание..." : "Создать"}
        </button>
      </div>
    </div>
  );
}
