CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    event_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    total_amount DECIMAL(20, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_pending_expires ON bookings(status, expires_at) WHERE status = 'PENDING';