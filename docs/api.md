## API

The data warehouse additionally exposes a couple of REST API's (GraphQL planned):

1. Dashboard API (`/dashboard`) - Exposes data for [`cic-dashboard`](https://github.com/grassrootseconomics/cic-dashboard). Most data is expected to be chart/table API specific and usually not human readable.
2. Public API (`/public`) - Exposes public (on-chain only/non-sensitive) data.
3. Internal Admin API - back office operations

Each API is domain separated i.e. separate SQL query files and router control.

### Pagination

Some API endpoints use a modified cursor pagination. This is to avaoid uncessary full db scans and clientside bloat when paginating data for tables.

The pagination expects the following query string:

- `?per_page=` (int) - No. of items to return. Has a hard limit of 100
- `?next=` (boolean) - If true, scrolls forward else backwards
- `?cursor=` (id:int) - Used with the forward query. If pagination forwards, pass in the id of the last result i.e. `results[length(results) - 1]`. If paginating backwards pass the id of the first element of the current result set i.e. `results[0]`.
