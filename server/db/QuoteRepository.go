package db

import (
	"botlord/models"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (db *BotlordDb) InsertQuote(quote models.Quote) (int, error) {
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

func (db *BotlordDb) DeleteQuote(id int) error {
	_, err := db.db.Exec("DELETE FROM quotes WHERE id = ?", id)
	if err != nil {
		fmt.Printf("error while deleting quote: %v\n", err)
		return err
	}
	return nil
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

func (b *BotlordDb) GetAllQuotes() ([]models.Quote, error) {
	rows, err := b.db.Query("SELECT id, text, lastchanged FROM quotes")
	if err != nil {
		fmt.Printf("error querying quotes: %v", err)
		return nil, err
	}
	defer rows.Close()

	quotes := []models.Quote{}
	for rows.Next() {
		var id int
		var text, lastchanged string
		if err := rows.Scan(&id, &text, &lastchanged); err != nil {
			return nil, fmt.Errorf("error scanning quote row: %w", err)
		}

		time, _ := time.Parse(time.RFC3339, lastchanged)

		quote := models.Quote{
			Id:          id,
			Text:        text,
			LastChanged: time,
		}
		quotes = append(quotes, quote)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("error iterating quote rows: %v", err)

		return nil, err
	}

	return quotes, nil
}
