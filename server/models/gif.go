package models

import (
	"fmt"
	"strings"
	"time"
)

type Gif struct {
	Id          int
	Url        string
	LastChanged time.Time
}

func NewGif(url string) *Gif {
	return &Gif{
		Url:        url,
		LastChanged: time.Now(),
	}
}

func PrintGifs(gifs []Gif) string {
	var result strings.Builder

	for _, q := range gifs {
		fmt.Fprintf(&result, "Id: %-5d Text: %s\n", q.Id, q.Url)
	}

	return result.String()
}

