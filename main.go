package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	sqlcdb "github.com/aitva/sqlc-parser-bug/db"
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
		if flags.MigrateDown {
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

		queries := sqlcdb.New(tx)

		if flags.Create {
			err := createMessage(queries)
			if err != nil {
				return fmt.Errorf("fail to create message: %v", err)
			}
		}

		if flags.List {
			err := listMessages(queries)
			if err != nil {
				return fmt.Errorf("fail to list messages: %v", err)
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
	MigrateDown bool
	List        bool
	Create      bool
}

func loadFlags() *flags {
	flags := &flags{}
	flag.BoolVar(&flags.MigrateUp, "up", false, "run migration up")
	flag.BoolVar(&flags.MigrateDown, "down", false, "run migration down")
	flag.BoolVar(&flags.List, "list", false, "run the list query")
	flag.BoolVar(&flags.Create, "create", false, "run the create query")
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
