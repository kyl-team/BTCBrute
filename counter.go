package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter int
	mu      sync.Mutex
)

func IncrementCounter() {
	mu.Lock()
	defer mu.Unlock()
	counter++
}

func ResetCounter() {
	mu.Lock()
	defer mu.Unlock()
	counter = 0
}

func GetCounter() int {
	mu.Lock()
	defer mu.Unlock()
	return counter
}

func counterDaemon() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C // Wait for the ticker to tick
		wallets := GetCounter()
		speed := wallets / 60

		sendMessage(fmt.Sprintf("Solved %d wallets. Speed is %d wallets/s", wallets, speed))
		ResetCounter()
	}
}
