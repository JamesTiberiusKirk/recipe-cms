package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type dbc struct {
	DB     *sqlx.DB
	schema goyesql.Queries
}

func newDBClient(dbUrl string) (*dbc, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	schema := goyesql.MustParseFile("./sql/schema.sql")

	return &dbc{
		DB:     db,
		schema: schema,
	}, nil
}

func (d *dbc) SchemaDown() error {
	sq, ok := d.schema["schema_down"]
	if !ok {
		return fmt.Errorf("schemanot not found")
	}

	tx := d.DB.MustBegin()

	tx.MustExec(sq.Query)
	tx.MustExec(`DROP TABLE IF EXISTS migrations;`)

	err := tx.Commit()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error("Rollback error")
			return err
		}
		return err
	}

	return nil
}

func (d *dbc) SchemaUp() error {
	sq, ok := d.schema["schema_up"]
	if !ok {
		return fmt.Errorf("schemanot not found")
	}

	tx := d.DB.MustBegin()

	tx.MustExec(sq.Query)
	tx.MustExec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id      SERIAL PRIMARY KEY,
			version INTEGER NOT NULL,
		);
		INSERT INTO migrations (id, version)
		VALUES (1, $1)
		ON CONFLICT (id)
		DO UPDATE SET version = EXCLUDED.version;
	`, countMigrations())

	err := tx.Commit()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error("Rollback error")
			return err
		}
		return err
	}

	return nil
}

func listFilesFilter(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func countMigrations() int {
	files, err := listFilesFilter("./sql/migrations", "*.sql")
	if err != nil {
		logrus.Errorf("Error opening migrations directory: %s", err.Error())
		panic(err)
	}

	return len(files)
}

func sortArray(arr []int) []int {
	for i := 0; i <= len(arr)-1; i++ {
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

func obfuscatePassword(connURL string) (string, error) {
	parsedURL, err := url.Parse(connURL)
	if err != nil {
		return "", err
	}

	username := parsedURL.User.Username()
	parsedURL.User = url.UserPassword(username, "xxxxxxx")

	return parsedURL.String(), nil
}

func applySchemaUp(dbc *dbc) {
	err := dbc.SchemaUp()
	if err != nil {
		logrus.Infof("Error applying schema up: %s", err.Error())
		panic(err)
	}
}

func applySchemaDown(dbc *dbc) {
	err := dbc.SchemaDown()
	if err != nil {
		logrus.Infof("Error applying schema down: %s", err.Error())
		panic(err)
	}
}

func applyMigration(dbc *dbc) {
	type row struct {
		ID      string `db:"id"`
		Version int    `db:"version"`
	}
	var r row
	err := dbc.DB.QueryRowx("SELECT * FROM migrations WHERE id = 1").StructScan(&r)
	if err != nil {
		logrus.Errorf("Error quering migrations table: %s", err.Error())
		panic(err)
	}
	logrus.Infof("Curent migration level: %d", r.Version)

	files, err := os.ReadDir("./sql/migrations")
	if err != nil {
		logrus.Errorf("Error opening migrations directory: %s", err.Error())
		panic(err)
	}

	var toApply []int

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		split := strings.Split(file.Name(), ".")
		if len(split) == 2 {
			continue
		}

		if split[1] != "sql" {
			continue
		}

		level, err := strconv.Atoi(split[0])
		if err != nil {
			logrus.Errorf("Could not parse migrations: %s", err.Error())
			panic(err)
		}

		if level > r.Version {
			toApply = append(toApply, level)
		}
	}

	if len(toApply) == 0 {
		logrus.Info("No new migrations")
		return
	}

	if len(toApply) > 1 {
		toApply = sortArray(toApply)
	}

	for _, l := range toApply {
		migration, err := os.ReadFile(fmt.Sprintf("./sql/migrations/%d.sql", l))
		if err != nil {
			logrus.Errorf("Could not read migration file %d: %s", l, err.Error())
			panic(err)
		}

		tx, err := dbc.DB.Begin()
		if err != nil {
			logrus.Errorf("Error begining transaction: %s", err.Error())
			panic(err)
		}

		_, err = tx.Exec(string(migration))
		if err != nil {
			logrus.Errorf("Error executing migration itself: %s", err.Error())
			panic(err)
		}

		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO migrations (id, version)
			VALUES (1, %d)
			ON CONFLICT (id)
			DO UPDATE SET version = EXCLUDED.version;
		`, l))
		if err != nil {
			logrus.Errorf("Error executing version upgrate in db transaction: %s", err.Error())
			panic(err)
		}

		err = tx.Commit()
		if err != nil {
			logrus.Error("failed to commit transaction")
			panic(err)
		}

		logrus.Infof("Applied migration: %d", l)
		logrus.Infof("Upgraded migration version number: %d", l)
	}
}

func connectDB() (string, *dbc) {
	dbUrl := os.Getenv("DB_URL")

	dbclient, err := newDBClient(dbUrl)
	if err != nil {
		panic(err)
	}

	url, err := obfuscatePassword(dbUrl)
	if err != nil {
		panic(err)
	}
	return url, dbclient
}

func main() {
	logrus.Info("------------------------------------------------------------")
	logrus.Info("MIGRATOR")

	action := flag.String("action", "", "[count-migrations|schema-up|schema-down|migrate]")
	flag.Parse()

	switch *action {
	case "count-migrations":
		logrus.Info("Available migrations")
		logrus.Info(countMigrations())
	case "schema-up":
		url, dbclient := connectDB()
		logrus.Infof("Applying schema up to db: %s", url)
		applySchemaUp(dbclient)
	case "schema-down":
		url, dbclient := connectDB()
		logrus.Infof("Applying schema down to db: %s", url)
		applySchemaDown(dbclient)
	case "migrate":
		url, dbclient := connectDB()
		logrus.Infof("Applying migration to db: %s", url)
		applyMigration(dbclient)
	default:
		flag.PrintDefaults()
	}
	logrus.Info("------------------------------------------------------------")
}
