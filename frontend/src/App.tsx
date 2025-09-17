import { OrderList } from "./components/Orders/OrderList"
import { ProductList } from "./components/Products/ProductList"


function App() {
  return (
    <>
      <div>
        <h1>Order Management</h1>
        <ProductList />
        <OrderList />
      </div>
    </>
  )
}

export default App
