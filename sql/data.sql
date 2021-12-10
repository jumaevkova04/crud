INSERT INTO products (name, price, qty)
VALUES ('Pizza', 200, 10),
    ('Burger', 150, 20),
    ('Free', 120, 15),
    ('Tea', 100, 50),
    ('Cola', 100, 50),
    ('Coffee', 100, 50);
-- 
--
INSERT INTO managers (name, salary, plan, boss_id, department)
VALUES ('Vasya', 100, 0, NULL, NULL),
    ('Petya', 80, 80, 1, 'boys'),
    ('Vanya', 60, 60, 2, 'boys'),
    ('Dasha', 90, 90, 1, 'girls'),
    ('Sasha', 70, 70, 4, 'girls'),
    ('Masha', 50, 50, 5, 'girls');
-- 
--
INSERT INTO customers (name, phone)
VALUES ('Zhenya', '+992000000001');
-- 
-- 
INSERT INTO sales (manager_id, customer_id)
VALUES (1, DEFAULT),
    (2, DEFAULT),
    (3, DEFAULT),
    (4, 1),
    (4, 1),
    (5, DEFAULT),
    (5, DEFAULT);
-- 
--
INSERT INTO sale_positions (sale_id, product_id, name, qty, price)
VALUES (1, 1, 'Pizza', 5, 200),
    (1, 2, 'Burger', 5, 200),
    (2, 3, 'Free', 10, 120),
    (3, 3, 'Free', 10, 120),
    (4, 6, 'Coffee', 20, 150),
    (5, 6, 'Coffee', 20, 150),
    (6, 6, 'Coffee', 20, 150),
    (7, 5, 'Cola', 10, 100);