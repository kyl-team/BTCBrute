package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func worker(token string, chatID string, id int, wg *sync.WaitGroup, mutex *sync.Mutex, outputFile string, btcAddresses map[string]bool) {
	defer wg.Done()

	for {
		privateKey, publicAddress, err := generateKeyAndAddress()
		consoleCounter.Add(1)
		reportCounter.Add(1)
		if err != nil {
			log.Printf("Worker %d: Failed to generate key and address: %s", id, err)
			continue
		}

		if _, exists := btcAddresses[publicAddress]; exists {
			fmt.Printf("Match Found! Privatekey: %s Publicaddress: %s\n", privateKey, publicAddress)

			mutex.Lock()
			file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Worker %d: Failed to open file: %s", id, err)
				mutex.Unlock()
				continue
			}

			logString := fmt.Sprintf("%s:%s\n", privateKey, publicAddress)
			sendMessage(token, chatID, fmt.Sprintf("Found new wallet: %s", logString))

			if _, err := file.WriteString(logString); err != nil {
				log.Printf("Worker %d: Failed to write to file: %s", id, err)
			}
			file.Close()
			mutex.Unlock()
		}
	}
}

func main() {
	outputFile := "data/output.txt"
	btcAddressesFile := "data/btc.txt"

	var threads int
	var token, chatID string
	flag.IntVar(&threads, "threads", 1, "Threads to run")
	flag.StringVar(&token, "token", "", "Telegram bot token")
	flag.StringVar(&chatID, "chatID", "", "Telegram chat ID")
	flag.Parse()

	sendMessage(token, chatID, "Loading address database")

	startTime := time.Now()

	btcAddresses, err := readAddresses(btcAddressesFile)
	if err != nil {
		log.Fatalf("Failed to read BTC addresses: %s", err)
	}

	sendMessage(token, chatID, fmt.Sprintf("Database loading took %.1fs. Loaded %d known wallets", time.Since(startTime).Seconds(), len(btcAddresses)))

	go reportDaemon(token, chatID)
	go consoleDaemon()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(token, chatID, i, &wg, &mutex, outputFile, btcAddresses)
	}

	wg.Wait()
}
