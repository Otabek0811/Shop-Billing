CREATE TABLE sale (
    id UUID PRIMARY KEY,
    branch_id UUID REFERENCES branch(id),
    assistant_id UUID REFERENCES staff(id),
    cashier_id UUID NOT NULL REFERENCES staff(id) ,
    price NUMERIC  NULL,
    payment_type VARCHAR NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'success',
    client_name VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL
);
