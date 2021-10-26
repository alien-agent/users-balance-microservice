INSERT INTO transaction (sender_id, recipient_id, amount, description, transaction_date)
VALUES (NULL, '11111111-3a7a-4d5e-8a6c-febc8c5b3f13', 2000, 'VISA TRANSFER', now()),
       (NULL, '22222222-3a7a-4d5e-8a6c-febc8c5b3f13', 3590, 'MASTERCARD TRANSFER', now() + interval '13s'),
       (NULL, '33333333-3a7a-4d5e-8a6c-febc8c5b3f13', 150, 'PROMO', now() + interval '25s'),
       ('22222222-3a7a-4d5e-8a6c-febc8c5b3f13', '11111111-3a7a-4d5e-8a6c-febc8c5b3f13', 1000, 'FOR DINNER', now() + interval '33s');

INSERT INTO deposit (owner_id, balance)
VALUES ('11111111-3a7a-4d5e-8a6c-febc8c5b3f13', 3000),
       ('22222222-3a7a-4d5e-8a6c-febc8c5b3f13', 2590),
       ('33333333-3a7a-4d5e-8a6c-febc8c5b3f13', 150);
