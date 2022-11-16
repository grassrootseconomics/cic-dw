-- name: get-password-hash
SELECT password_hash FROM staff WHERE username = $1;

-- name: pin-status
SELECT phone_number, failed_pin_attempts,
CASE STATUS
	WHEN 1 THEN 'PENDING'
    WHEN 2 THEN 'ACTIVE'
    WHEN 3 THEN 'LOCKED'
    WHEN 4 THEN 'RESET' END AS account_status
FROM cic_ussd.account WHERE
failed_pin_attempts > 0 OR STATUS = 4;

-- name: phone-2-address
SELECT blockchain_address FROM users WHERE phone_number = $1;

-- name: address-2-phone
SELECT phone_number FROM users WHERE blockchain_address = $1;

-- name: account-latest-transactions
-- Returns the first page of a users latest transactions
SELECT transactions.id, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
FROM transactions
INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
INNER JOIN tokens ON transactions.token_address = tokens.token_address
WHERE users.phone_number = $1 ORDER BY transactions.id DESC LIMIT $2;

-- name: account-latest-transactions-next
-- Returns the next page based on a cursor
SELECT transactions.id, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
FROM transactions
INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
INNER JOIN tokens ON transactions.token_address = tokens.token_address
WHERE users.phone_number = $1 AND transactions.id < $2 ORDER BY transactions.id DESC LIMIT $3;

-- name: account-latest-transactions-previous
-- Returns the previous page based on cursor
SELECT * FROM (
    SELECT transactions.id, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
    (SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
    (SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
    FROM transactions
    INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
    INNER JOIN tokens ON transactions.token_address = tokens.token_address
    WHERE users.phone_number = $1 AND transactions.id > $2 ORDER BY transactions.id ASC LIMIT $3
) AS previous_page ORDER BY id DESC;

-- name: account-latest-transactions-by-token
-- Returns the first page of a users latest transactions, filter by token
SELECT transactions.id, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
FROM transactions
INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
INNER JOIN tokens ON transactions.token_address = tokens.token_address
WHERE users.phone_number = $1 AND tokens.token_symbol = $2 ORDER BY transactions.id DESC LIMIT $3;

-- name: account-latest-transactions-by-token-next
-- Returns the next page based on a cursor, filter by token
SELECT transactions.id, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
FROM transactions
INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
INNER JOIN tokens ON transactions.token_address = tokens.token_address
WHERE users.phone_number = $1 AND tokens.token_symbol = $2 AND transactions.id < $3 ORDER BY transactions.id DESC LIMIT $4;

-- name: account-latest-transactions-by-token-previous
-- Returns the previous page based on cursor, filter by token
SELECT * FROM (
    SELECT transactions.id, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
    (SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
    (SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
    FROM transactions
    INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
    INNER JOIN tokens ON transactions.token_address = tokens.token_address
    WHERE users.phone_number = $1 AND tokens.token_symbol = $2 AND transactions.id > $3 ORDER BY transactions.id ASC LIMIT $4
) AS previous_page ORDER BY id DESC;

-- name: account-latest-transactions-by-archived-token
-- Returns the first page of a users latest transactions, filter by token
SELECT transactions.id, transactions.date_block, transactions.tx_hash, archived_tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_addres) as recipient_identifier
FROM transactions
INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
INNER JOIN archived_tokens ON transactions.token_address = archived_tokens.token_address
WHERE users.phone_number = $1 AND archived_tokens.token_address = $2 ORDER BY transactions.id DESC LIMIT $3;

-- name: account-latest-transactions-by-archived-token-next
-- Returns the next page based on a cursor, filter by token
SELECT transactions.id, transactions.date_block, transactions.tx_hash, archived_tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
(SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
FROM transactions
INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
INNER JOIN archived_tokens ON transactions.token_address = archived_tokens.token_address
WHERE users.phone_number = $1 AND archived_tokens.token_address = $2 AND transactions.id < $3 ORDER BY transactions.id DESC LIMIT $4;

-- name: account-latest-transactions-by-archived-token-previous
-- Returns the previous page based on cursor, filter by token
SELECT * FROM (
    SELECT transactions.id, transactions.date_block, transactions.tx_hash, archived_tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success,
    (SELECT phone_number FROM users WHERE blockchain_address = transactions.sender_address) as sender_identifier,
    (SELECT phone_number FROM users WHERE blockchain_address = transactions.recipient_address) as recipient_identifier
    FROM transactions
    INNER JOIN users ON ((transactions.sender_address = users.blockchain_address) OR (transactions.recipient_address = users.blockchain_address))
    INNER JOIN archived_tokens ON transactions.token_address = archived_tokens.token_address
    WHERE users.phone_number = $1 AND archived_tokens.token_address = $2 AND transactions.id > $3 ORDER BY transactions.id ASC LIMIT $4
) AS previous_page ORDER BY id DESC;
