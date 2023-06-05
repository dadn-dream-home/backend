server {
    port = 50051
}

database {
    driver = "sqlite3_with_hook"
    connection_string = "./db.sqlite3"
    migrations_path = "./migrations"
}

mqtt {
    brokers = ["tcp://localhost:1883"]
}