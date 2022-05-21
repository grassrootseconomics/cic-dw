CREATE TABLE IF NOT EXISTS sys_accounts (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sys_address VARCHAR(40) NOT NULL,
    address_description VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS sys_address_idx ON sys_accounts USING hash(sys_address);

INSERT INTO sys_accounts(sys_address, address_description) VALUES
('bbb4a93c8dcd82465b73a143f00fed4af7492a27', 'sarafu sink address'),
('cd9fd1e71f684cfb30fa34831ed7ed59f6f77469', 'sarafu faucet'),
('b8830b647c01433f9492f315ddbfdc35cb3be6a6', 'ge community fund'),
('ca5da01b6dac771c8f3625aa1a8931e7dac41832', 'ge token deployer'),
('59a5e2faf8163fe24ca006a221dd0f34c5e0cb41', 'sarafu migrator'),
('289defd53e2d96f05ba29ebbebd9806c94d04cb6', 'sohail deployer');

---- create above / drop below ----

DROP TABLE IF EXISTS sys_accounts;