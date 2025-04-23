package main

import (
	//"log"
	"os"
	"os/signal"
	"syscall"

	"botlord/bot"
)


func main() {
	b := bot.NewBot()	
	b.Start()
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
