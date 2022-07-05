-- This change makes token management manual

ALTER TABLE tokens DROP CONSTRAINT tokens_token_address_key;