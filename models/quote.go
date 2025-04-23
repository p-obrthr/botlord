package models

import(
	"time"
)

type Quote struct {
	Id int
	Text string
	LastChanged time.Time
}

func NewQuote(text string) *Quote {
	return &Quote{
		Text:        text,
		LastChanged: time.Now(),
	}
}
