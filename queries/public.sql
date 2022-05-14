-- name: all-known-tokens
-- Looks up all known tokens from the transactions records
SELECT DISTINCT tokens.token_symbol, tokens.token_address FROM transactions
INNER JOIN tokens on transactions.token_address = tokens.token_address
WHERE transactions.sender_address = $1
OR transactions.recipient_address = $1;
