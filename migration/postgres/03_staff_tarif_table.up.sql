CREATE TABLE staff_tarif (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    type_tarif VARCHAR NOT NULL,
    AmountForCash   NUMERIC DEFAULT 0,
    AmountForCard   NUMERIC DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL
);
