## Syncers

Syncers are responsible for pulling in data from multiple sources into this central data warehouse. There are primarily 3 types of syncers:

1. On-chain syncer - Pulls in data from smart contracts storages, events e.t.c.
2. DB syncer - Sync from other Postgres databases through [`postgres_fdw`](https://www.postgresql.org/docs/current/postgres-fdw.html) which may rely on other syncers on their end. 
3. API syncer - Pull in data from an exposed endpoint.

Syncers are implemented as periodic jobs and rely on cursor based pagination across all types of syncers to efficiently update the data. In some cases where the data is regularly updated, we either listen for events (if the remote is capable of sending events) or spawn periodic jobs for every entry.

### Syncers implemented

#### Cache syncer

Relies on the `cic_cache` db which is backed by cic-cache-tracker. Pulls in all on-chain transfer transactions.

#### Ussd syncer

Relies on the `cic_ussd` db. Pulls in relevant ussd account registration details and doubles up as an alternative to the on-chain account registry index.

#### Token syncer

On-chain syncer. Pulls in all registered CIC tokens from the token index.

#### Meta syncer (WIP)

API syncer. Pulls in all relevant user metadata from the CRDT backed meta API.