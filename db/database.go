package db

import (
	"database/sql"
	"path/filepath"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"botlord/models"
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

func(db *BotlordDb) Insert(quote models.Quote) (int, error) {
	now := time.Now()
	timestamp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	res, err := db.db.Exec("INSERT INTO quotes (text, lastchanged) VALUES (?,?);", quote.Text, timestamp)
	if err != nil {
		fmt.Printf("err while inserting quote: db.db.exec %v", err)
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *BotlordDb) GetRandomQuoteText() (*string, error) {
	row := db.db.QueryRow("SELECT text FROM quotes ORDER BY RANDOM() LIMIT 1;")

	var quote string
	err := row.Scan(&quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

