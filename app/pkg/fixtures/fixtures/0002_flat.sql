INSERT INTO flat (house_id, flat_number, price, rooms, status, created_at, updated_at)
VALUES
    (1, 1, 400, 3, 'created', now(), now()),
    (1, 2, 800, 4, 'approved', now(), now()),
    (1, 3, 100, 1, 'on moderation', now(), now()),
    (3, 1, 200, 1, 'created', now(), now());
