package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var counter atomic.Uint64

func counterDaemon() {
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		wallets := counter.Load()
		speed := wallets / 3600 / 4

		sendMessage(fmt.Sprintf("Solved %d wallets. Speed is %d wallets/s", wallets, speed))
		counter.Store(0)
	}
}
