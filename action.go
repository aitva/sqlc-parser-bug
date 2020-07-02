package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aitva/sqlc-parser-bug/db"
	"github.com/google/uuid"
)

const migrationFolder = "./schema"

func migrateUp(tx *sql.Tx) error {
	return filepath.Walk(migrationFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk: %v", err)
		}
		if info.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %v: %v", path, err)
		}

		_, err = tx.Exec(string(data))
		if err != nil {
			return fmt.Errorf("exec %v: %v", path, err)
		}

		return nil
	})
}

func migrateDrop(tx *sql.Tx, dbuser string) error {
	_, err := tx.Exec("DROP OWNED BY " + dbuser)
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}
	return nil
}

func createMessage(q *db.Queries) error {
	externalID, _ := uuid.NewRandom()
	_, err := q.CreateMessage(context.Background(), db.CreateMessageParams{
		Content:    "Hello sqlc!",
		ExternalID: []uuid.UUID{externalID},
	})
	if err != nil {
		return err
	}
	return nil
}

func createCounter(q *db.Queries) error {
	_, err := q.CreateCounter(context.Background(), []int64{1})
	if err != nil {
		return err
	}
	return nil
}
