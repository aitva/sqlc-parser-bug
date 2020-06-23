package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	flags := loadFlags()

	info, err := loadDBInfo()
	if err != nil {
		fmt.Printf("fail to load DB infos: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", info.String())
	if err != nil {
		fmt.Printf("fail to open DB connection: %v\n", err)
		os.Exit(1)
	}

	err = runInTx(db, func(tx *sql.Tx) error {
		if flags.MigrateDrop {
			err := migrateDrop(tx, info.User)
			if err != nil {
				return fmt.Errorf("fail to drop database: %v", err)
			}
		}

		if flags.MigrateUp {
			err := migrateUp(tx)
			if err != nil {
				return fmt.Errorf("fail to migrate up: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

type flags struct {
	MigrateUp   bool
	MigrateDrop bool
}

func loadFlags() *flags {
	flags := &flags{}
	flag.BoolVar(&flags.MigrateUp, "up", false, "run migration up")
	flag.BoolVar(&flags.MigrateDrop, "drop", false, "drop the database")
	flag.Parse()
	return flags
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin: %v", err)
	}

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit: %v", err)
	}

	return nil
}
