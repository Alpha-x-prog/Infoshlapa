package telegram

import (
	"database/sql"
	"log"
	"newsAPI/db"
	"os/exec"
	"strings"
)

type ChannelRequest struct {
	ChannelURL string `json:"channel_url" binding:"required"`
}

// AddChannel добавляет новый канал в базу данных и запускает Python скрипт для его обработки
func AddChannel(dbConn *sql.DB, userID int, req ChannelRequest) error {
	// Извлекаем username из URL
	channelUsername := strings.TrimPrefix(req.ChannelURL, "https://t.me/")
	if channelUsername == req.ChannelURL {
		channelUsername = strings.TrimPrefix(req.ChannelURL, "t.me/")
	}

	// Запускаем Python скрипт для получения информации о канале
	cmd := exec.Command("python", "telegram/scripts/get_channel_info.py", req.ChannelURL)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error getting channel info: %v", err)
		return err
	}

	// Получаем название канала из вывода скрипта
	channelTitle := strings.TrimSpace(string(output))
	if channelTitle == "" {
		channelTitle = channelUsername // Используем username как fallback
	}

	// Сначала сохраняем канал в таблицу user_channels
	err = db.AddUserChannel(dbConn, userID, req.ChannelURL, channelUsername, channelTitle)
	if err != nil {
		log.Printf("Error saving channel to user_channels: %v", err)
		return err
	}

	// Проверяем, существует ли канал в telegram_channels
	var exists bool
	err = dbConn.QueryRow("SELECT EXISTS(SELECT 1 FROM telegram_channels WHERE channel_username = ?)", channelUsername).Scan(&exists)
	if err != nil {
		log.Printf("Error checking channel existence: %v", err)
		return err
	}

	// Если канала нет в telegram_channels, добавляем его
	if !exists {
		_, err = dbConn.Exec("INSERT INTO telegram_channels (channel_username, channel_title, last_message_id, is_active) VALUES (?, ?, ?, ?)",
			channelUsername, channelTitle, 0, true)
		if err != nil {
			log.Printf("Error adding channel to telegram_channels: %v", err)
			return err
		}

		// Запускаем Python скрипт для обработки канала
		cmd = exec.Command("python", "telegram/scripts/add_channel.py", req.ChannelURL)
		if err := cmd.Run(); err != nil {
			log.Printf("Error running Python script: %v", err)
			// В случае ошибки удаляем канал из обеих таблиц
			_ = db.RemoveUserChannel(dbConn, userID, req.ChannelURL)
			_, _ = dbConn.Exec("DELETE FROM telegram_channels WHERE channel_username = ?", channelUsername)
			return err
		}
	}

	return nil
}

// RemoveChannel удаляет канал из базы данных
func RemoveChannel(dbConn *sql.DB, userID int, channelURL string) error {
	return db.RemoveUserChannel(dbConn, userID, channelURL)
}

// GetUserChannels получает все каналы пользователя
func GetUserChannels(dbConn *sql.DB, userID int) ([]map[string]string, error) {
	return db.GetUserChannels(dbConn, userID)
}
