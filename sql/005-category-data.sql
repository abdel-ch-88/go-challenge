-- Insert 3 categories
INSERT INTO categories (code, name) VALUES
('CAT001', 'Clothing'),
('CAT002', 'Shoes'),
('CAT003', 'Accessories');

-- Set category for each product using category code to look up category_id
--      and match products by thier code

-- Category 1: 3 products
UPDATE products
SET category_id = (SELECT id FROM categories WHERE code = 'CAT001')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

-- Category 2: 2 products
UPDATE products
SET category_id = (SELECT id FROM categories WHERE code = 'CAT002')
WHERE code IN ('PROD002', 'PROD006');

-- Category 3: 3 products
UPDATE products
SET category_id = (SELECT id FROM categories WHERE code = 'CAT003')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');
