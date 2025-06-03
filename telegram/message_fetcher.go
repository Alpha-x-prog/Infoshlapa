package telegram

import (
	"log"
	"os/exec"
	"time"
)

// StartMessageFetcher запускает фоновый процесс получения сообщений
func StartMessageFetcher() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				cmd := exec.Command("python", "telegram/scripts/fetch_messages.py")
				if err := cmd.Run(); err != nil {
					log.Printf("Error running fetch_messages.py: %v", err)
				}
			}
		}
	}()
	log.Println("Message fetcher started")
}
