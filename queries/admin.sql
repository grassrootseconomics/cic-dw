-- name: get-password-hash
SELECT password_hash from staff WHERE username = $1;