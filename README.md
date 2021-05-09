# Use with docker
## docker-compose
### Configure bot
Copy -dist files and set content

    $ cat db.env
    POSTGRES_PASSWORD=koalasmuertos
    POSTGRES_USER=diosteama
    POSTGRES_DB=diosteama

    $ cat bot.env
    TELEGRAM_BOT_TOKEN=secret:token
    DIOSTEAMA_DB_URL=postgres://diosteama:koalasmuertos@db/diosteama

### Launch bot
Execute database

    docker-compose up -d db

Import dump

    zcat ~/Downloads/quotes.2020-06-22T10_00+00_00.sql.gz |docker exec -i diosteama_db_1 psql -U diosteama

Execute bot

    docker-compose up -d bot

