# Nubayrah
Ebook library and metadata manager

# How to Run
Run the following command to have docker build the golang server and also instantiate postgresqls server.

`docker compose up -d --build`

The default path for database storage is `./.data/`. Please change this if needed.

# Current API

`GET /books` Returns JSON of all items in database.

`GET /books/{bookID}` Returns specified json item.

`POST /books` Sends a json body and creates the entry inside of DB.

`DELETE /books/{bookID}` Deletes entry by id in database.