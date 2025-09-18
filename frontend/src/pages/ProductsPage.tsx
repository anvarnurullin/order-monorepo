import { ProductList } from "../components/Products/ProductList";
import './ProductsPage.scss';

export default function ProductsPage() {
  return (
    <div className="products-page">
      <h1 className="products-page__title">Каталог товаров</h1>
      <ProductList />
    </div>
  );
}