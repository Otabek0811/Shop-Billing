CREATE TABLE staff_transaction(
    id UUID PRIMARY KEY,
    transaction_type VARCHAR DEFAULT 'Topup',
    description VARCHAR ,
    price NUMERIC NOT NULL,
    sale_id UUID REFERENCES sale(id),
    source_type VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL
);