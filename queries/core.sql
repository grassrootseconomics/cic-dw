-- syncer

-- name: ussd-syncer
-- This db transaction will auto scroll the cic_ussd remote adding values as per the limit and auto-updating the cursor
-- The blockchain_address is used as a cursor to retrieve the corresponding id since the id is not guaranteed to be sequential
WITH current_ussd_cursor AS (
    SELECT id FROM cic_ussd.account WHERE blockchain_address = (SELECT cursor_pos FROM cursors WHERE id = 1)
)

INSERT INTO users (phone_number, blockchain_address, date_registered, failed_pin_attempts, ussd_account_status)
SELECT cic_ussd.account.phone_number, cic_ussd.account.blockchain_address, cic_ussd.account.created, cic_ussd.account.failed_pin_attempts, cic_ussd.account.status
FROM cic_ussd.account WHERE cic_ussd.account.id > (SELECT id FROM current_ussd_cursor) ORDER BY cic_ussd.account.id ASC LIMIT 300;

UPDATE cursors SET cursor_pos = (SELECT blockchain_address FROM users ORDER BY id DESC LIMIT 1) WHERE cursors.id = 1;

-- name: cache-syncer
-- This db transaction will auto scroll the cic_cache remote adding values as per the limit and auto-updating the cursor
-- The tx_hash is used as the cursor to retrieve the corresponding id since the id is not guaranteed to be sequential
WITH current_cache_cursor AS (
    SELECT id FROM cic_cache.tx WHERE LOWER(tx_hash) = (SELECT cursor_pos FROM cursors WHERE id = 2)
)

INSERT INTO transactions (tx_hash, block_number, tx_index, token_address, sender_address, recipient_address, tx_value, date_block, tx_type, success)
SELECT cic_cache.tx.tx_hash, cic_cache.tx.block_number, cic_cache.tx.tx_index, LOWER(cic_cache.tx.source_token), LOWER(cic_cache.tx.sender), LOWER(cic_cache.tx.recipient), cic_cache.tx.from_value, cic_cache.tx.date_block, concat(cic_cache.tag.domain, '_', cic_cache.tag.value) AS tx_type, cic_cache.tx.success
FROM cic_cache.tx INNER JOIN cic_cache.tag_tx_link ON cic_cache.tx.id = cic_cache.tag_tx_link.tx_id INNER JOIN cic_cache.tag ON cic_cache.tag_tx_link.tag_id = cic_cache.tag.id
WHERE cic_cache.tx.id > (SELECT id FROM current_cache_cursor) ORDER BY cic_cache.tx.id ASC LIMIT 300;

UPDATE cursors SET cursor_pos = (SELECT tx_hash FROM transactions ORDER BY id DESC LIMIT 1) WHERE cursors.id = 2;

-- name: cursor-pos
-- Generic cursor query
SELECT cursor_pos from cursors WHERE id = $1;

-- name: insert-token-data
-- Insert new token
INSERT INTO tokens (token_address, token_name, token_symbol, token_decimals) VALUES
    (LOWER($1), $2, $3, $4);

-- name: update-cursor
-- Generic cursor update
UPDATE cursors SET cursor_pos = $1 WHERE cursors.id = $2;