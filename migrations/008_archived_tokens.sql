CREATE TABLE IF NOT EXISTS archived_tokens (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_address VARCHAR(40) NOT NULL UNIQUE,
    token_decimals INT NOT NULL,
    token_name VARCHAR(16) NOT NULL,
    token_symbol VARCHAR(10) NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS archived_tokens;