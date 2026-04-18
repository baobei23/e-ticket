TRUNCATE TABLE events RESTART IDENTITY CASCADE;

INSERT INTO events (name, description, location, start_time, end_time) VALUES
('Konser Coldplay', 'Music of the Spheres', 'GBK',
 NOW() + INTERVAL '1 month', NOW() + INTERVAL '1 month 4 hours'),
('Java Jazz Festival', 'International Jazz Music Festival', 'JIExpo Kemayoran',
 NOW() + INTERVAL '2 months', NOW() + INTERVAL '2 months 8 hours'),
('We The Fest', 'Music, Arts & Fashion Festival', 'GBK Sports Complex',
 NOW() + INTERVAL '3 months', NOW() + INTERVAL '3 months 10 hours');

-- Categories
INSERT INTO seat_categories (event_id, name, price, color, sort_order) VALUES
(1, 'VIP',     5000000, '#FFD700', 1),
(1, 'Reguler', 2500000, '#4F46E5', 2),
(2, 'VIP',     2500000, '#FFD700', 1),
(2, 'Reguler', 1200000, '#4F46E5', 2),
(3, 'VIP',     2800000, '#FFD700', 1),
(3, 'Reguler', 1500000, '#4F46E5', 2);

-- Generate seats: tiap event 2 kategori
-- Event 1: VIP row A-B (masing-masing 10 kursi) + Reguler row C-F (masing-masing 20)
DO $$
DECLARE
    e RECORD;
    c RECORD;
    r CHAR;
    n INT;
    vip_rows CHAR[] := ARRAY['A','B'];
    reg_rows CHAR[] := ARRAY['C','D','E','F'];
BEGIN
    FOR e IN SELECT id FROM events LOOP
        -- VIP seats
        SELECT id INTO c FROM seat_categories WHERE event_id = e.id AND name = 'VIP';
        FOREACH r IN ARRAY vip_rows LOOP
            FOR n IN 1..10 LOOP
                INSERT INTO seats (event_id, category_id, row_label, seat_number)
                VALUES (e.id, c.id, r, n);
            END LOOP;
        END LOOP;

        -- Reguler seats
        SELECT id INTO c FROM seat_categories WHERE event_id = e.id AND name = 'Reguler';
        FOREACH r IN ARRAY reg_rows LOOP
            FOR n IN 1..20 LOOP
                INSERT INTO seats (event_id, category_id, row_label, seat_number)
                VALUES (e.id, c.id, r, n);
            END LOOP;
        END LOOP;
    END LOOP;
END $$