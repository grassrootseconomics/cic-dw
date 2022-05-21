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