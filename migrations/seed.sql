-- Products
INSERT INTO products (id, name, description, images, category, created_at, updated_at) VALUES
(
    'product_001',
    'Classic White T-Shirt',
    'A comfortable everyday white t-shirt made from 100% cotton.',
    ARRAY['https://placehold.co/600x400?text=White+TShirt'],
    'Clothing',
    NOW(),
    NOW()
),
(
    'product_002',
    'Black Hoodie',
    'A warm and stylish black hoodie perfect for casual wear.',
    ARRAY['https://placehold.co/600x400?text=Black+Hoodie'],
    'Clothing',
    NOW(),
    NOW()
),
(
    'product_003',
    'Blue Jeans',
    'Classic slim fit blue jeans for everyday wear.',
    ARRAY['https://placehold.co/600x400?text=Blue+Jeans'],
    'Clothing',
    NOW(),
    NOW()
),
(
    'product_004',
    'Running Shoes',
    'Lightweight and comfortable running shoes.',
    ARRAY['https://placehold.co/600x400?text=Running+Shoes'],
    'Footwear',
    NOW(),
    NOW()
);

-- SKUs
INSERT INTO skus (id, product_id, sku_code, price, stock, attributes, created_at, updated_at) VALUES
-- White T-Shirt SKUs
('sku_001', 'product_001', 'WHITE-TSHIRT-S-WHITE', 200, 10, '{"size": "S", "color": "white"}', NOW(), NOW()),
('sku_002', 'product_001', 'WHITE-TSHIRT-M-WHITE', 200, 15, '{"size": "M", "color": "white"}', NOW(), NOW()),
('sku_003', 'product_001', 'WHITE-TSHIRT-L-WHITE', 200, 8, '{"size": "L", "color": "white"}', NOW(), NOW()),
('sku_004', 'product_001', 'WHITE-TSHIRT-S-BLACK', 220, 5, '{"size": "S", "color": "black"}', NOW(), NOW()),
('sku_005', 'product_001', 'WHITE-TSHIRT-M-BLACK', 220, 12, '{"size": "M", "color": "black"}', NOW(), NOW()),

-- Black Hoodie SKUs
('sku_006', 'product_002', 'BLACK-HOODIE-S-BLACK', 450, 6, '{"size": "S", "color": "black"}', NOW(), NOW()),
('sku_007', 'product_002', 'BLACK-HOODIE-M-BLACK', 450, 10, '{"size": "M", "color": "black"}', NOW(), NOW()),
('sku_008', 'product_002', 'BLACK-HOODIE-L-BLACK', 450, 4, '{"size": "L", "color": "black"}', NOW(), NOW()),
('sku_009', 'product_002', 'BLACK-HOODIE-M-GRAY',  430, 8, '{"size": "M", "color": "gray"}',  NOW(), NOW()),

-- Blue Jeans SKUs
('sku_010', 'product_003', 'BLUE-JEANS-30-BLUE', 580, 7, '{"size": "30", "color": "blue"}',  NOW(), NOW()),
('sku_011', 'product_003', 'BLUE-JEANS-32-BLUE', 580, 9, '{"size": "32", "color": "blue"}',  NOW(), NOW()),
('sku_012', 'product_003', 'BLUE-JEANS-34-BLUE', 580, 3, '{"size": "34", "color": "blue"}',  NOW(), NOW()),
('sku_013', 'product_003', 'BLUE-JEANS-32-BLACK', 600, 5, '{"size": "32", "color": "black"}', NOW(), NOW()),

-- Running Shoes SKUs
('sku_014', 'product_004', 'RUNNING-SHOES-40-WHITE', 780, 4, '{"size": "40", "color": "white"}', NOW(), NOW()),
('sku_015', 'product_004', 'RUNNING-SHOES-42-WHITE', 780, 6, '{"size": "42", "color": "white"}', NOW(), NOW()),
('sku_016', 'product_004', 'RUNNING-SHOES-44-WHITE', 780, 2, '{"size": "44", "color": "white"}', NOW(), NOW()),
('sku_017', 'product_004', 'RUNNING-SHOES-42-BLACK', 800, 5, '{"size": "42", "color": "black"}', NOW(), NOW());