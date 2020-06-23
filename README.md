# sqlc-parser-bug

This repository reproduce a problem with the new
[sqlc](github.com/kyleconroy/sqlc) parser.

The folder `schema` contains two migrations:

```sql
-- schema/0001_messages.sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    content TEXT NOT NULL,
    external_id UUID
);
```

```sql
-- schema/0002_update_messages.sql
ALTER TABLE messages ALTER external_id TYPE UUID ARRAY USING ARRAY[external_id];
```

After running the migrations the type of the `external_id` field is `UUID[]`
but sqlc continue to treat it as `UUID`.

## How to

To run this application you need a local Postgres database. You can configure
the application using environment variables: `DB_HOST, DB_NAME, DB_PORT,
DB_USER, DB_PASS`.

1. configure the environment variables: `export DB_HOST=localhost DB_USER=test
   DB_PASS=test`
1. setup the database: `go run . -up`
1. create a record: `go run . -create`

The create record command will fail with the error: `fail to create message: pq:
malformed array literal: "7de80066-106d-4b7e-8c74-8cfe6202f5b7"`.

## Fix

Run sqlc with the old parser: `SQLC_EXPERIMENTAL_PARSER=off sqlc generate`.

The generated code will correctly detect the type of `external_id` and a build
error will catch the problem: `./action.go:51:3: cannot use externalID (type
uuid.UUID) as type []uuid.UUID in field value`.
