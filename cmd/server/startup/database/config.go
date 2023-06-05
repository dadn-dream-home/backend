package database

type Config struct {
	ConnectionString string `hcl:"connection_string"`
	MigrationsPath   string `hcl:"migrations_path"`
	Driver           string `hcl:"driver"`
}
