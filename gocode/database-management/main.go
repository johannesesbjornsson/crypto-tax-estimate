package main

import (
    "github.com/go-gorp/gorp/v3"
		"github.com/lib/pq"
		log "github.com/sirupsen/logrus"
		"os"
)


func main() {

}

func RunMigrations(db *sql.DB, migrationsDir string) error {
		db, err := sql.Open("postgres", "user=personal password=password dbname=cryptotax sslmode=disable")
		if err != nil {
		    log.Fatal(err)
		}

    files, err := ioutil.ReadDir(migrationsDir)
    if err != nil {
        return fmt.Errorf("reading migration dir: %w", err)
    }

    return nil
}