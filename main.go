package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

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

	if flags.HasMigrateDown {
		// TODO: drop all tables
		_ = db
	}

	if flags.HasMigrateUp {
		// TODO: migrate up
		_ = db
	}
}

type flags struct {
	HasMigrateUp   bool
	HasMigrateDown bool
}

func loadFlags() *flags {
	flags := &flags{}
	flag.BoolVar(&flags.HasMigrateUp, "up", false, "run migration up")
	flag.BoolVar(&flags.HasMigrateDown, "down", false, "run migration down")
	flag.Parse()
	return flags
}

type dbInfo struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

func (info *dbInfo) String() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		info.Host, info.Port, info.Name, info.User, info.Pass)
}

func loadDBInfo() (*dbInfo, error) {
	info := &dbInfo{
		Host: os.Getenv("DB_HOST"),
		Port: 5432,
		Name: os.Getenv("DB_NAME"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
	}
	if info.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if info.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if info.Name == "" {
		info.Name = info.User
	}
	if tmp := os.Getenv("DB_PORT"); tmp != "" {
		var err error
		info.Port, err = strconv.Atoi(tmp)
		if err != nil {
			return nil, fmt.Errorf("fail to parse DB_PORT: %v", err)
		}
	}
	return info, nil
}
