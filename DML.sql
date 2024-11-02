INSERT INTO customer (name, phone_number, address)
VALUES
('John Doe', '555-1234', '123 Elm St'),
('Jane Smith', '555-5678', '456 Oak St'),
('Michael Johnson', '555-8765', '789 Pine St'),
('Emily Davis', '555-4321', '321 Maple Ave'),
('Robert Brown', '555-1111', '654 Cedar St');

INSERT INTO employee (name, phone_number, address)
VALUES
('Alice Williams', '555-2222', '987 Willow St'),
('David Harris', '555-3333', '123 Birch St'),
('Sophia Martinez', '555-4444', '456 Redwood St'),
('James Wilson', '555-5555', '789 Palm St'),
('Olivia Garcia', '555-6666', '321 Cypress Ave');

INSERT INTO product (product_name, unit, price)
VALUES
('Shampoo', 'bottle', 10000),
('Soap', 'bar', 5000),
('Toothpaste', 'tube', 15000),
('Conditioner', 'bottle', 12000),
('Body Lotion', 'bottle', 25000);

INSERT INTO transaction (customer_id, employee_id, bill_date, entry_date, finish_date) 
VALUES 
(1, 1, '01-10-2024', '01-10-2024', '05-10-2024'),
(2, 2, '02-10-2024', '02-10-2024', '06-10-2024'),
(3, 3, '03-10-2024', '03-10-2024', '07-10-2024'),
(4, 4, '04-10-2024', '04-10-2024', '08-10-2024'),
(5, 5, '05-10-2024', '05-10-2024', '09-10-2024');

INSERT INTO transaction_detail (transaction_id, product_id, product_price, qty)
VALUES
(1, 1, 10000, 2),
(2, 2, 5000, 5),
(3, 3, 15000, 3),
(4, 4, 12000, 4),
(5, 5, 25000, 1);
