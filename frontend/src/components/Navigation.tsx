import { Link, useLocation } from "react-router-dom";
import './Navigation.scss';

export default function Navigation() {
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
    <nav className="navigation">
      <div className="navigation__container">
        <Link 
          to="/" 
          className={`navigation__link ${isActive('/') ? 'navigation__link--active' : ''}`}
        >
          Главная
        </Link>
        <Link 
          to="/products" 
          className={`navigation__link ${isActive('/products') ? 'navigation__link--active' : ''}`}
        >
          Товары
        </Link>
        <Link 
          to="/orders" 
          className={`navigation__link ${isActive('/orders') ? 'navigation__link--active' : ''}`}
        >
          Заказы
        </Link>
      </div>
    </nav>
  );
}