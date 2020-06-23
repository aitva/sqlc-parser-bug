package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
