-- Drop static columns on user table
ALTER TABLE users DROP COLUMN failed_pin_attempts, DROP COLUMN ussd_account_status;

-- Staff dashboard auth
CREATE TABLE IF NOT EXISTS staff (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(16) NOT NULL UNIQUE,
    password_hash VARCHAR(76) NOT NULL
)

---- create above / drop below ----

DROP TABLE IF EXISTS staff;