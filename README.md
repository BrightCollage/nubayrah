# Nubayrah
Ebook library and metadata manager

# How to Run

## Docker Compose
Run the following command to have docker build the golang server and also instantiate postgresqls server.

`docker compose up -d --build`

The default path for database storage is `./.data/`. Please change this if needed.

## go run

Command to run API server:
`go run ./cmd/api`

Command to run HTML server:
`go run ./cmd/html`

# Current API

`GET /books` Returns JSON of all items in database.

`GET /books/{id}` Returns specified json item.

`POST /books` Sends a json body and creates the entry inside of DB.

`DELETE /books/{id}` Deletes entry by id in database.

`GET /books/{id}/cover` Returns image of specified item.

# Client

We're using ReactJS + Vite as our framework. The project lives in `/client`.

## Build

To build the client, run:

`npm --prefix client run build`

The destination for the build will be the `./static` directory, which will be served by the `go` httpserver.

Or you can cd into the `client` directory first and then run:

`npm run build`