-- name: all-known-tokens
-- Looks up all known tokens from the transactions records
SELECT DISTINCT tokens.token_symbol, tokens.token_address FROM transactions
INNER JOIN tokens on transactions.token_address = tokens.token_address
WHERE transactions.sender_address = $1
OR transactions.recipient_address = $1;

-- Bidirectional cursor paginators
-- name: list-tokens-fwd
SELECT tokens.id, tokens.token_address, tokens.token_name, tokens.token_symbol FROM tokens
WHERE tokens.id > $1 ORDER BY tokens.id ASC LIMIT $2;

-- name: list-tokens-bkwd
SELECT tokens.id, tokens.token_address, tokens.token_name, tokens.token_symbol FROM tokens
WHERE tokens.id < $1 ORDER BY tokens.id ASC LIMIT $2;

-- name: tokens-count
-- Return total record count from individual i= tables/views
SELECT COUNT(*) FROM tokens;


--name: unique-token-holders
-- Returns the unique token holders based on seen transactions
WITH unique_holders AS (
	SELECT sender_address AS holding_address FROM transactions
  	WHERE token_address = $1
  	UNION
  	SELECT recipient_address AS holding_address FROM transactions
  	WHERE token_address = $1
),
exclude AS (
    SELECT sys_address FROM sys_accounts WHERE sys_address IS NOT NULL
)

SELECT COUNT(holding_address) FROM unique_holders
WHERE holding_address NOT IN (SELECT sys_address FROM exclude);

--name: all-time-token-transactions-count
-- Returns transactions of individual tokens
WITH exclude AS (
    SELECT sys_address FROM sys_accounts WHERE sys_address IS NOT NULL
)

SELECT COUNT(*) FROM transactions
WHERE token_address = $1
AND transactions.sender_address NOT IN (SELECT sys_address FROM exclude)
AND transactions.recipient_address NOT IN (SELECT sys_address FROM exclude)
AND transactions.success = true;


--name: latest-token-transactions
-- Returns latest token transactions, with curosr forward query and limit
SELECT transactions.id, transactions.block_number, transactions.date_block, transactions.tx_hash, tokens.token_symbol, transactions.sender_address, transactions.recipient_address, transactions.tx_value, transactions.success FROM transactions
INNER JOIN tokens ON transactions.token_address = tokens.token_address
WHERE transactions.token_address = $1 AND transactions.id < $2 ORDER BY transactions.id DESC LIMIT $3;