package main

import (
	//"github.com/go-gorp/gorp/v3"
	"database/sql"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
  "path/filepath"
  "sort"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	db, err := sql.Open("postgres", "user=personal password=password dbname=cryptotax sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	RunMigrations(db, "migrations")

}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return err
	}
	var sqlFiles []os.FileInfo
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
		log.Infof("Found migration file: %s", file.Name())
	}

	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		path := filepath.Join(migrationsDir, file.Name())
		query, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		log.Infof("🟡 Running migration: %s", file.Name())
		if _, err := db.Exec(string(query)); err != nil {
			log.Warnf("❌ Failed migration: %s. Error: %v", file.Name(), err)
		}
		log.Infof("✅ Completed migration: %s", file.Name())
	}

	return nil
}
