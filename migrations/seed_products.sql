INSERT INTO products (name, sku, price, qty_available) VALUES
('T-shirt', 'TS-001', 19.99, 100),
('Mug', 'MG-001', 9.50, 200)
ON CONFLICT DO NOTHING;
