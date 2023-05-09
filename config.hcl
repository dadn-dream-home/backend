server {
    port = 50051
}

database {
    connection_string = "./db.sqlite3"
    migrations_path = "./migrations"
}

mqtt {
    brokers = ["tcp://localhost:1883"]
}