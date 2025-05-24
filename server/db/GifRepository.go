package db

import (
	"botlord/models"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (db *BotlordDb) InsertGif(gif models.Gif) (int, error) {
	now := time.Now()
	timestamp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	res, err := db.db.Exec("INSERT INTO gifs (url, lastchanged) VALUES (?,?);", gif.Url, timestamp)
	if err != nil {
		fmt.Printf("err while inserting gif: db.db.exec %v", err)
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *BotlordDb) DeleteGif(id int) error {
	_, err := db.db.Exec("DELETE FROM gifs WHERE id = ?", id)
	if err != nil {
		fmt.Printf("error while deleting gif: %v\n", err)
		return err
	}
	return nil
}

func (db *BotlordDb) GetRandomGif() (*string, error) {
	row := db.db.QueryRow("SELECT url FROM gifs ORDER BY RANDOM() LIMIT 1;")

	var gifUrl string
	err := row.Scan(&gifUrl)
	if err != nil {
		return nil, err
	}
	return &gifUrl, nil
}

func (b *BotlordDb) GetAllGifs() ([]models.Gif, error) {
	rows, err := b.db.Query("SELECT id, url, lastchanged FROM gifs")
	if err != nil {
		fmt.Printf("error querying quotes: %v", err)
		return nil, err
	}
	defer rows.Close()

	gifs := []models.Gif{}
	for rows.Next() {
		var id int
		var url, lastchanged string
		if err := rows.Scan(&id, &url, &lastchanged); err != nil {
			return nil, fmt.Errorf("error scanning gif row: %w", err)
		}

		time, _ := time.Parse(time.RFC3339, lastchanged)

		gif := models.Gif{
			Id:          id,
			Url:         url,
			LastChanged: time,
		}
		gifs = append(gifs, gif)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("error iterating gif rows: %v", err)

		return nil, err
	}

	return gifs, nil
}
