CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(255) PRIMARY KEY, -- UUID string
    booking_id VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    amount DECIMAL(20, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'idr',
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    stripe_id VARCHAR(255),
    payment_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);