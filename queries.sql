-- syncers

-- name: ussd-syncer
-- This db transaction will auto scroll the cic_ussd remote adding values as per the limit and auto-updating the cursor
-- The blockchain_address is used as a cursor to retrieve the corresponding id since the id is not guaranteed to be sequential
WITH current_ussd_cursor AS (
    SELECT id FROM cic_ussd.account WHERE blockchain_address = (SELECT cursor_pos FROM cursors WHERE id = 1)
)

INSERT INTO users (phone_number, blockchain_address, date_registered)
SELECT phone_number, blockchain_address, created
FROM cic_ussd.account WHERE id > (SELECT id FROM current_ussd_cursor) ORDER BY id ASC LIMIT 10;

UPDATE cursors SET cursor_pos = (SELECT blockchain_address FROM users ORDER BY id DESC LIMIT 1) WHERE cursors.id = 1;

-- name: cache-syncer
-- This db transaction will auto scroll the cic_cache remote adding values as per the limit and auto-updating the cursor
-- The tx_hash is used as the cursor to retrieve the corresponding id since the id is not guaranteed to be sequential
WITH current_cache_cursor AS (
    SELECT id FROM cic_cache.tx WHERE tx_hash = (SELECT cursor_pos FROM cursors WHERE id = 2)
)

INSERT INTO transactions (tx_hash, block_number, tx_index, token_address, sender_address, recipient_address, tx_value, date_block, tx_type)
SELECT tx.tx_hash, tx.block_number, tx.tx_index, tx.source_token, tx.sender, tx.recipient, tx.from_value, tx.date_block, concat(tag.domain, '_', tag.value) AS tx_type
FROM cic_cache.tx INNER JOIN cic_cache.tag_tx_link ON tx.id = cic_cache.tag_tx_link.tx_id INNER JOIN cic_cache.tag ON cic_cache.tag_tx_link.tag_id = cic_cache.tag.id
WHERE tx.success = true AND tx.id > (SELECT id FROM current_cache_cursor) ORDER BY tx.id ASC LIMIT 10;

UPDATE cursors SET cursor_pos = (SELECT tx_hash FROM tx ORDER BY id DESC LIMIT 1) WHERE cursors.id = 2;