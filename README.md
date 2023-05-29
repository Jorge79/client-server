This is a personal project fo studies purposes. It's a simple web application that allows users to see the dollar value in reals and save the value in a database.

# About
This project was made using Go, GORM, SQLite and Docker.

## How to run
It's necessary to have Docker/Docker Compose installed in your machine. After that, run the following commands:
```
docker compose up -d
```

Change the directory to the server folder and run the command:
```
go run .
```
Open a new terminal and change the directory to the client folder and run the command:
```
curl localhost:8080/cotacao
```

To see the database, enter in the container:
```
docker ps
docker exec -it <container_id> bash

sqlite3 db.sqlite3
```

