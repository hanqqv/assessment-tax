CREATE TYPE allowance_type AS ENUM ('personal');

CREATE TABLE IF NOT EXISTS deductions_setting (
    id SERIAL PRIMARY KEY,
    allowance_type allowance_type NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO deductions_setting (allowance_type, amount) VALUES 
('personal', 60000.00);