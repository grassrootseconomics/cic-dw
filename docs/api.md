## API

The data warehouse additionally exposes a couple of REST API's (GraphQL planned):

1. Dashboard API (`/dashboard`) - Exposes data for [`cic-dashboard`](https://github.com/grassrootseconomics/cic-dashboard). Most data is expected to be chart/table API specific and usually not human readable.  
2. Public API (`/public`) - Exposes public (on-chain only/non-sensitive) data.
3. Internal API (planned)

Each API is domain separated i.e. separate SQL query files and router control. 