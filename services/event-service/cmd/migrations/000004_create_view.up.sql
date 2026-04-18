CREATE OR REPLACE VIEW v_event_availability AS
SELECT
    e.id AS event_id,
    COUNT(s.id) FILTER (WHERE s.status = 'AVAILABLE') AS available_seats,
    COUNT(s.id) AS total_seats,
    MIN(sc.price) AS min_price,
    MAX(sc.price) AS max_price
FROM events e
LEFT JOIN seats s ON s.event_id = e.id
LEFT JOIN seat_categories sc ON sc.event_id = e.id
GROUP BY e.id;
