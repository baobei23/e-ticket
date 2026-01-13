CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    total_seats INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- Seed Data (Optional)
INSERT INTO events (name, description, location, start_time, end_time, total_seats, available_seats, price)
VALUES ('Konser Coldplay', 'Music of the Spheres', 'GBK', NOW() + INTERVAL '1 month', NOW() + INTERVAL '1 month 4 hours', 50000, 50000, 3500000);