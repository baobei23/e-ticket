CREATE TABLE IF NOT EXISTS seat_categories (
  id BIGSERIAL PRIMARY KEY,
  event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  price BIGINT NOT NULL,
  color VARCHAR(255),
  sort_order INT NOT NULL DEFAULT 0,
  UNIQUE (event_id, name)
);

CREATE INDEX idx_seat_categories_event ON seat_categories(event_id);