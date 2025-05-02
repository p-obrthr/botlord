package main

import (
	"os"
	"os/signal"
	"syscall"

	"botlord/bot"
	"botlord/api"
)

func main() {
	var wr *api.Wrapper

	if os.Getenv("ENABLE_API") == "1" {
		wr = api.NewWrapper()
	}

	b := bot.NewBot()
	b.Start()

	if wr != nil {
		wr.Bot = b
		wr.Running = true
		wr.StartHTTPServer(":8080")
	}
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if wr != nil {
		wr.StopBot()
	} else {
		b.Stop()
	}
}
