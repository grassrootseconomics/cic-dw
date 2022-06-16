CREATE UNIQUE INDEX index_id ON transactions USING btree (ID DESC);

---- create above / drop below ----

DROP INDEX index_id;