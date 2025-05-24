package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type BotlordDb struct {
	db *sql.DB
}

const dbDir string = "db"
const file string = "botlord.db"

const createDb string = `
	CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER NOT NULL PRIMARY KEY,
		text TEXT NOT NULL,
		lastchanged TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS gifs (
		id INTEGER NOT NULL PRIMARY KEY,
		url TEXT NOT NULL,
		lastchanged TEXT NOT NULL
	);

`

func InitDb() (*BotlordDb, error) {
	dbPath := filepath.Join(dbDir, file)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createDb); err != nil {
		fmt.Println("createDb err")
		return nil, err
	}

	botlordDb := &BotlordDb{db: db}

	return botlordDb, nil
}



