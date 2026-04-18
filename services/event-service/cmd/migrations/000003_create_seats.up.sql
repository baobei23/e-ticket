CREATE TABLE IF NOT EXISTS seats (
  id BIGSERIAL PRIMARY KEY,
  event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  category_id BIGINT NOT NULL REFERENCES seat_categories(id) ON DELETE CASCADE,
  row_label VARCHAR(8) NOT NULL,      
  seat_number INT NOT NULL,
  label VARCHAR(16) GENERATED ALWAYS AS (row_label || seat_number::text) STORED,
  status VARCHAR(16) NOT NULL DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE', 'HELD', 'SOLD')),
  held_by_booking_id UUID,
  held_until TIMESTAMP,
  UNIQUE (event_id, row_label, seat_number)
);

CREATE INDEX idx_seats_event_status ON seats(event_id, status);
CREATE INDEX idx_seat_held_booking ON seats(held_by_booking_id) WHERE status = 'HELD';
CREATE INDEX idx_seats_held_until ON seats(held_until) WHERE status = 'HELD';