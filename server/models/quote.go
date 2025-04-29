package models

import (
	"fmt"
	"strings"
	"time"
)

type Quote struct {
	Id          int
	Text        string
	LastChanged time.Time
}

func NewQuote(text string) *Quote {
	return &Quote{
		Text:        text,
		LastChanged: time.Now(),
	}
}

func PrintQuotes(quotes []Quote) string {
	var result strings.Builder

	for _, q := range quotes {
		fmt.Fprintf(&result, "Id: %-5d Text: %s\n", q.Id, q.Text)
	}

	return result.String()
}

