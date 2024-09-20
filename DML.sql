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
(1, 1, 20240910, '2024-09-10 08:30:00', '2024-09-12 10:30:00'),
(2, 2, 20240911, '2024-09-11 09:00:00', '2024-09-13 12:00:00'),
(3, 3, 20240912, '2024-09-12 10:00:00', '2024-09-14 11:00:00'),
(4, 4, 20240913, '2024-09-13 11:30:00', '2024-09-15 14:00:00'),
(5, 5, 20240914, '2024-09-14 12:00:00', '2024-09-16 15:30:00');

INSERT INTO transaction_detail (transaction_id, product_id, product_price, qty)
VALUES
(1, 1, 10000, 2),
(2, 2, 5000, 5),
(3, 3, 15000, 3),
(4, 4, 12000, 4),
(5, 5, 25000, 1);
