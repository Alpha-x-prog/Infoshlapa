package handlers

import (
	"database/sql"
	"net/http"
	"newsAPI/db"

	"github.com/gin-gonic/gin"
)

var dbConn *sql.DB

// InitDB initializes the database connection for handlers
func InitDB(db *sql.DB) {
	dbConn = db
}

// GetPublicChannels получает список публичных каналов
func GetPublicChannels(c *gin.Context) {
	if dbConn == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	rows, err := dbConn.Query(`
		SELECT channel_username, channel_title
		FROM telegram_channels
		WHERE channel_username IN ('priem_mirea', 'mirea_esports', 'mireaprofkom')
		ORDER BY channel_title
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
		return
	}
	defer rows.Close()

	var channels []map[string]string
	for rows.Next() {
		channel := make(map[string]string)
		var username, title string
		if err := rows.Scan(&username, &title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan channel data"})
			return
		}
		channel["username"] = username
		channel["title"] = title
		channels = append(channels, channel)
	}

	c.JSON(http.StatusOK, channels)
}

// GetPublicChannelMessages получает сообщения из публичных каналов
func GetPublicChannelMessages(c *gin.Context) {
	if dbConn == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	messages, err := db.GetPublicChannelMessages(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}
