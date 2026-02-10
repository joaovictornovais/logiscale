CREATE TABLE IF NOT EXISTS locations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id UUID NOT NULL,
    lat DECIMAL(9,6) NOT NULL,
    lng DECIMAL(9,6) NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL,
    
    FOREIGN KEY (driver_id) REFERENCES drivers(id) ON DELETE CASCADE
);