import React, { useEffect, useState } from "react";
import type { Product } from "../types/product";
import { fetchProducts } from "../api/products";
import { ImageUpload } from "./ImageUpload";
import "./ProductList.scss";


export const ProductList: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);

  useEffect(() => {
    fetchProducts().then(setProducts);
  }, []);

  const handleImageUploaded = (productId: number, imageUrl: string) => {
    setProducts(prev => 
      prev.map(p => 
        p.id === productId ? { ...p, imageUrl } : p
      )
    );
  };

  return (
    <div className="product-list">
      <div className="product-list__grid">
        {products.sort((a, b) => a.id - b.id).map((p) => (
          <div key={p.id} className="product-list__card">
            <h3 className="product-list__title">{p.name}</h3>
            <p className="product-list__info">SKU: {p.sku}</p>
            <p className="product-list__info">Цена: ${p.price}</p>
            <p className="product-list__info">Количество: {p.qtyAvailable}</p>
            
            {p.imageUrl ? (
              <div>
                <img 
                  src={p.imageUrl} 
                  alt={p.name}
                  className="product-list__image"
                  onError={(e) => {
                    console.error('Failed to load image:', p.imageUrl);
                    e.currentTarget.style.display = 'none';
                  }}
                  onLoad={() => {
                    console.log('Image loaded successfully:', p.imageUrl);
                  }}
                />
              </div>
            ) : (
              <div className="product-list__placeholder">
                Нет изображения
              </div>
            )}
            
            <ImageUpload 
              productId={p.id} 
              onImageUploaded={(imageUrl) => handleImageUploaded(p.id, imageUrl)}
            />
          </div>
        ))}
      </div>
    </div>
  );
};
