-- tx table
CREATE TABLE IF NOT EXISTS transactions (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tx_hash VARCHAR(64) NOT NULL,
    block_number INT NOT NULL,
    tx_index INT NOT NULL,
    token_address VARCHAR(40) NOT NULL,
    sender_address VARCHAR(40) NOT NULL,
    recipient_address VARCHAR(40) NOT NULL,
    tx_value BIGINT NOT NULL,
    tx_type VARCHAR(16) NOT NULL,
    date_block TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS token_idx ON transactions USING hash(token_address);
CREATE INDEX IF NOT EXISTS sender_idx ON transactions USING hash(sender_address);
CREATE INDEX IF NOT EXISTS recipient_idx ON transactions USING hash(recipient_address);

-- tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_address VARCHAR(40) NOT NULL UNIQUE,
    token_decimals INT NOT NULL,
    token_name VARCHAR(10) NOT NULL,
    token_symbol VARCHAR(10) NOT NULL
);

-- users table
CREATE TABLE IF NOT EXISTS users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    phone_number VARCHAR(16) NOT NULL,
    blockchain_address VARCHAR(40) NOT NULL,
    date_registered TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS phone_number_idx ON users USING hash(phone_number);
CREATE INDEX IF NOT EXISTS sender_idx ON users USING hash(blockchain_address);

-- trigram extension for location and product search
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

-- meta table
CREATE TABLE IF NOT EXISTS meta (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    preferred_language VARCHAR(16),
    gender VARCHAR(10),
    age INT,
    given_name VARCHAR(32),
    family_name VARCHAR(32),
    products TEXT [],
    location_name VARCHAR(32),
    location_latitude FLOAT,
    location_longitude FLOAT,
    tags TEXT[]
);

CREATE INDEX IF NOT EXISTS tags ON meta USING gin(tags);
CREATE INDEX IF NOT EXISTS location_name_idx ON meta USING gin(location_name);
CREATE INDEX IF NOT EXISTS products_idx ON meta USING gin(location_name);
CREATE INDEX IF NOT EXISTS meta_filter_idx ON meta(gender, preferred_language, age);

-- cursors table (for internal syncing)
CREATE TABLE IF NOT EXISTS cursors (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cursor_pos VARCHAR(64) NOT NULL
);

-- bootstrap first users row
INSERT INTO users (phone_number, blockchain_address, date_registered)
SELECT phone_number, blockchain_address, created
FROM cic_ussd.account WHERE id = 1;

-- idx 1 = cic_ussd cursor
INSERT INTO cursors (cursor_pos)
SELECT blockchain_address FROM users ORDER BY id DESC LIMIT 1;

-- bootstrap first tx row
INSERT INTO transactions (tx_hash, block_number, tx_index, token_address, sender_address, recipient_address, tx_value, date_block, tx_type)
SELECT tx.tx_hash, tx.block_number, tx.tx_index, tx.source_token, tx.sender, tx.recipient, tx.from_value, tx.date_block, concat(tag.domain, '_', tag.value) AS tx_type
FROM cic_cache.tx
INNER JOIN cic_cache.tag_tx_link ON tx.id = cic_cache.tag_tx_link.tx_id
INNER JOIN cic_cache.tag ON cic_cache.tag_tx_link.tag_id = cic_cache.tag.id
WHERE tx.success = true AND tx.id = 1;

-- idx 2 = cic_cache cursor
INSERT INTO cursors (cursor_pos)
SELECT tx_hash FROM transactions ORDER BY id DESC LIMIT 1;
