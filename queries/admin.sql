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

--name: phone-2-address
SELECT blockchain_address FROM users WHERE phone_number = $1;

--name: address-2-phone
SELECT phone_number FROM users WHERE blockchain_address = $1;