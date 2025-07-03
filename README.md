# Echoes

Twitter clone

## Stack

- Go 1.24.3

## Setup

1. Create a `.env` file. it should have following variables:
    ```env
    ENV=
    DB_URL=
    AUTH_SECRET=
    POLKA_KEY=
    ```
2. Run a Postgres. You can use docker compose, see `compose.yaml` file
3. Run Goose UP migrations manually or via make, see `Makefile` file
    ```
    make migrate
    make migrate-down
    ```

## Build

```
go build
```

## Run

```
./echoes
```
or
```
go run .
```

