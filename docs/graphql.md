# GraphQL

## Endpoints

The API provides a GraphQL query and sandbox endpoint.

| METHOD | ENDPOINT                       | DESCRIPTION                    | 
|--------|--------------------------------|--------------------------------|
| GET    | /q/playground                  | Interactive GraphQL playground |
| POST   | /q/query                       | GraphQL query endpoint         |

## Queries

The following queries are supported:

| QUERY                       | DESCRIPTION                      | 
|-----------------------------|----------------------------------|
| show(name: "name")          | Returns a show and its episodes  |
| episode(guid: "guid")       | Returns a single episode, identified by its GUID |
| recent(max: int)            | Returns up to 'max' recently updated shows, based on their build date, in descending order |

