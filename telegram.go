package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func sendMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botToken, chatID, url.QueryEscape(message))
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
