
CREATE TABLE   branch(
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    company_id UUID REFERENCES company(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at  TIMESTAMP DEFAULT NULL
);