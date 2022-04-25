-- fdw extension for remote link
CREATE EXTENSION IF NOT EXISTS postgres_fdw;

-- cic_ussd
CREATE SCHEMA IF NOT EXISTS cic_ussd;
CREATE SERVER IF NOT EXISTS cic_ussd_remote FOREIGN DATA WRAPPER postgres_fdw OPTIONS
    (host '{{.remote_db_host }}', port '{{.remote_db_port }}', dbname 'cic_ussd');
CREATE USER MAPPING IF NOT EXISTS FOR postgres SERVER cic_ussd_remote OPTIONS
    (user '{{.remote_db_user }}', password '{{.remote_db_password }}');
IMPORT FOREIGN SCHEMA public LIMIT TO (account) FROM SERVER cic_ussd_remote INTO cic_ussd;

-- cic_cache
CREATE SCHEMA IF NOT EXISTS cic_cache;
CREATE SERVER IF NOT EXISTS cic_cache_remote FOREIGN DATA WRAPPER postgres_fdw OPTIONS
(host '{{.remote_db_host }}', port '{{.remote_db_port }}', dbname 'cic_cache');
CREATE USER MAPPING IF NOT EXISTS FOR postgres SERVER cic_cache_remote OPTIONS
(user '{{.remote_db_user }}', password '{{.remote_db_password }}');
IMPORT FOREIGN SCHEMA public LIMIT TO (tag, tag_tx_link, tx) FROM SERVER cic_cache_remote INTO cic_cache;
