package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var user, password, host, port, databaseName, migrationsPath, migrationsTable string

	flag.StringVar(&user, "user", "postgres", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.StringVar(&port, "port", "5432", "port")
	flag.StringVar(&databaseName, "database", "", "database name")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if password == "" {
		panic("password is required")
	}
	if databaseName == "" {
		panic("database is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?x-migrations-table=%s&sslmode=disable", user, password, host, port, databaseName, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}
}
