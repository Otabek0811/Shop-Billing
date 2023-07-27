CREATE TABLE staff(
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    staff_type VARCHAR NOT NULL,
    balance NUMERIC DEFAULT 0,
    tarif_id UUID REFERENCES staff_tarif(id),
    branch_id UUID REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL
);
