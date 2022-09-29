-- name: new-user-registrations
-- This query generates a date range and left joins the users table to include days with no registrations
-- Produces x, y results for displaying on a line chart
WITH date_range AS (
    SELECT day::date FROM generate_series($1, $2, INTERVAL '1 day') day
)

SELECT date_range.day AS x, COUNT(users.id) AS y
FROM date_range
LEFT JOIN users ON date_range.day = CAST(users.date_registered AS date)
GROUP BY date_range.day
ORDER BY date_range.day
LIMIT 730;

-- name: new-user-registrations-cmr
-- This is a patch to support CMR dashboard
-- This query generates a date range and left joins the users table to include days with no registrations
-- Produces x, y results for displaying on a line chart
WITH date_range AS (
    SELECT day::date FROM generate_series($1, $2, INTERVAL '1 day') day
)

SELECT date_range.day AS x, COUNT(users.id) AS y
FROM date_range
LEFT JOIN users ON date_range.day = CAST(users.date_registered AS date)
WHERE users.phone_number LIKE '237%'
GROUP BY date_range.day
ORDER BY date_range.day
LIMIT 730;

-- name: transactions-count
-- This query generates a date range and left joins the transactions table to include days with no transactions
-- Produces x, y results for displaying on a line chart
WITH date_range AS (
    SELECT day::date FROM generate_series($1, $2, INTERVAL '1 day') day
),
exclude AS (
    SELECT sys_address FROM sys_accounts WHERE sys_address IS NOT NULL
)

SELECT date_range.day AS x, COUNT(transactions.id) AS y
FROM date_range
LEFT JOIN transactions ON date_range.day = CAST(transactions.date_block AS date)
AND transactions.sender_address NOT IN (SELECT sys_address FROM exclude) AND transactions.recipient_address NOT IN (SELECT sys_address FROM exclude)
AND transactions.success = true
GROUP BY date_range.day
ORDER BY date_range.day
LIMIT 730;

-- name: transactions-count-cmr
-- This is a patch to support CMR dashboard
-- This query generates a date range and left joins the transactions table to include days with no transactions
-- Produces x, y results for displaying on a line chart
WITH date_range AS (
    SELECT day::date FROM generate_series($1, $2, INTERVAL '1 day') day
),
exclude AS (
    SELECT sys_address FROM sys_accounts WHERE sys_address IS NOT NULL
),
cmr_tokens AS (
  	SELECT token_address from tokens WHERE token_symbol IN ('MBIP', 'MUN')
)
SELECT date_range.day AS x, COUNT(transactions.id) AS y
FROM date_range
LEFT JOIN transactions ON date_range.day = CAST(transactions.date_block AS date)
AND transactions.sender_address NOT IN (SELECT sys_address FROM exclude) AND transactions.recipient_address NOT IN (SELECT sys_address FROM exclude)
AND transactions.success = true
AND transactions.token_address IN (SELECT token_address FROM cmr_tokens)
GROUP BY date_range.day
ORDER BY date_range.day
LIMIT 730;

-- name: token-transactions-count
-- This query gets transactions for a specific token for a given date range
WITH date_range AS (
    SELECT day::date FROM generate_series($1, $2, INTERVAL '1 day') day
),
exclude AS (
    SELECT sys_address FROM sys_accounts WHERE sys_address IS NOT NULL
)

SELECT date_range.day AS x, COUNT(transactions.id) AS y
FROM date_range
LEFT JOIN transactions ON date_range.day = CAST(transactions.date_block AS date)
AND transactions.sender_address NOT IN (SELECT sys_address FROM exclude) AND transactions.recipient_address NOT IN (SELECT sys_address FROM exclude)
AND transactions.token_address = $3
AND transactions.success = true
GROUP BY date_range.day
ORDER BY date_range.day
LIMIT 730;

--name: token-volume
-- This query rteurns daily token volume
-- Assumes erc20 token is 6 decimals
WITH date_range AS (
    SELECT day::date FROM generate_series($1, $2, INTERVAL '1 day') day
),
exclude AS (
    SELECT sys_address FROM sys_accounts WHERE sys_address IS NOT NULL
)

SELECT date_range.day AS x, COALESCE(SUM(transactions.tx_value / 1000000), 0) AS y
FROM date_range
LEFT JOIN transactions ON date_range.day = CAST(transactions.date_block AS date)
AND transactions.sender_address NOT IN (SELECT sys_address FROM exclude) AND transactions.recipient_address NOT IN (SELECT sys_address FROM exclude)
AND transactions.token_address = $3
AND transactions.success = true
GROUP BY date_range.day
ORDER BY date_range.day
LIMIT 730;

--name: latest-token-transactions
-- Returns latest token transactions, with curosr forward query and limit
SELECT id, block_number, date_block, tx_hash, sender_address, recipient_address, tx_value, success FROM transactions
WHERE token_address = $1 AND date_block > TIMESTAMP 'yesterday' ORDER BY id DESC;
