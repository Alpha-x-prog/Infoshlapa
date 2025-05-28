package telegram

import (
	"log"
	"os/exec"
)

// AddChannel запускает Python скрипт для добавления нового канала
func AddChannel(channelURL string) error {
	cmd := exec.Command("python", "telegram/scripts/add_channel.py", channelURL)
	if err := cmd.Run(); err != nil {
		log.Printf("Error adding channel: %v", err)
		return err
	}
	return nil
}
