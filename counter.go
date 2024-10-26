package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var reportCounter atomic.Uint64
var consoleCounter atomic.Uint64

func reportDaemon(botToken string, chatID string) {
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		wallets := reportCounter.Load()
		speed := wallets / 3600 / 4

		sendMessage(botToken, chatID, fmt.Sprintf("Solved %d wallets. Speed is %d wallets/s", wallets, speed))
		reportCounter.Store(0)
	}
}

func consoleDaemon() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		println(fmt.Sprintf("Solved %d wallets in 30 seconds", consoleCounter.Load()))
		consoleCounter.Store(0)
	}
}
