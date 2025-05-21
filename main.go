package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"newsAPI/db" // Замените на актуальный путь к пакету db
	"newsAPI/gemini"
	_ "newsAPI/gemini"
	"newsAPI/parser"
	_ "newsAPI/parser"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем JWT ключ из переменных окружения
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY не найден в переменных окружения")
	}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJWT(userID int) (string, error) {
	// Устанавливаем время истечения токена (например, 24 часа)
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	// Получаем API ключ из переменных окружения
	apiKey := os.Getenv("NEWSDATA_API_KEY")
	if apiKey == "" {
		log.Fatal("NEWSDATA_API_KEY не найден в переменных окружения")
	}

	// Инициализация подключения к БД
	database, err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка инициализации БД:", err)
		return
	}
	defer database.Close()

	categories := []string{"top", "health", "politics", "sports", "business", "science", "food"}

	//// Запускаем горутину для каждого типа категории
	for _, category := range categories {
		go startNewsFetcher(apiKey, category, database)
	}

	r := gin.Default()

	// Раздача статических файлов
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/news", func(c *gin.Context) {
		getNews(c, database)
	})

	//помощник
	r.POST("/ask", func(c *gin.Context) {
		geminiASK(c)
	})

	// Public routes
	r.POST("/register", func(c *gin.Context) {
		RegisterHandler(c, database)
	})

	r.POST("/login", func(c *gin.Context) {
		LoginHandler(c, database)
	})

	// Protected routes (require JWT)
	protected := r.Group("/protected")
	protected.Use(JWTAuthMiddleware())
	{
		// Example protected route
		protected.GET("/profile", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "User ID not found in context"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome to your profile!", "user_id": userID})
		})
	}

	r.Run(":8080")
}

type NewsArticle struct {
	ArticleID   string   `json:"article_id"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Keywords    []string `json:"keywords"`
	Creator     []string `json:"creator"`
	VideoURL    string   `json:"video_url"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	PubDate     string   `json:"publishedAt"`
	ImageURL    string   `json:"urlToImage"`
	SourceID    string   `json:"source_id"`
	SourceName  string   `json:"source_name"`
	SourceURL   string   `json:"url"`
	Language    string   `json:"language"`
	Country     string   `json:"country"`
	Tags        string   `json:"tags"`
	Sentiment   string   `json:"sentiment"`
}

func getNews(c *gin.Context, database *sql.DB) {
	limit := 15
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	category := c.DefaultQuery("category", "")

	validCategories := []string{"top", "sports", "technology", "business", "science", "entertainment", "health", "world", "politics", "environment", "food"}
	validCategory := false

	for _, valid := range validCategories {
		if category == valid {
			validCategory = true
			break
		}
	}

	log.Printf("Получение новостей с лимитом %d, смещением %d, категорией: %s", limit, offset, category)

	var rows *sql.Rows
	var err error

	if validCategory {
		rows, err = database.Query(
			"SELECT article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment "+
				"FROM news WHERE category LIKE ? ORDER BY pub_date DESC LIMIT ? OFFSET ?",
			"%"+category+"%", limit, offset,
		)
	} else {
		rows, err = database.Query(
			"SELECT article_id, title, link, keywords, creator, video_url, description, content, pub_date, image_url, source_id, source_name, source_url, language, country, category, sentiment "+
				"FROM news ORDER BY pub_date DESC LIMIT ? OFFSET ?",
			limit, offset,
		)
	}

	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось выполнить запрос"})
		return
	}
	defer rows.Close()

	var news []NewsArticle

	for rows.Next() {
		var n NewsArticle
		var keywordsStr, creatorStr, countryStr, categoryStr string

		err := rows.Scan(
			&n.ArticleID, &n.Title, &n.Link, &keywordsStr, &creatorStr, &n.VideoURL,
			&n.Description, &n.Content, &n.PubDate, &n.ImageURL, &n.SourceID, &n.SourceName, &n.SourceURL,
			&n.Language, &countryStr, &categoryStr, &n.Sentiment,
		)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
			return
		}

		// Разбиваем строки с разделителями в срезы
		n.Keywords = strings.Split(keywordsStr, ",")
		n.Creator = strings.Split(creatorStr, ",")
		n.Country = strings.TrimSpace(strings.Split(countryStr, ",")[0])
		n.Tags = strings.TrimSpace(strings.Split(categoryStr, ",")[0])

		news = append(news, n)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при обработке строк результата: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных"})
		return
	}

	c.JSON(http.StatusOK, news)
}

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Content string `json:"content"`
}

func geminiASK(c *gin.Context) {
	var req Request

	// Декодируем JSON из тела запроса в структуру Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования запроса"})
		return
	}
	userQuery := req.Prompt
	// Получаем ответ от функции geminiResponse

	responseContent := gemini.GeminiResponse("Напиши кратко ответ на вопрос: " + userQuery)
	// Отправляем JSON-ответ с полученным ответом
	c.JSON(http.StatusOK, Response{Content: responseContent})
}

func startNewsFetcher(apiKey, category string, database *sql.DB) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	apiURL := fmt.Sprintf("https://newsdata.io/api/1/latest?apikey=%%s&category=%s&language=ru&country=ru", category)

	for {
		log.Printf("Запуск парсинга для категории: %s", category)
		err := parser.ParseAndSaveNews(apiURL, apiKey, database)
		if err != nil {
			log.Printf("Ошибка при парсинге новостей (%s): %v", category, err)
		}
		<-ticker.C
	}
}

// --- User Registration and Login Handlers ---
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(c *gin.Context, db *sql.DB) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ?", req.Email, req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Server error"})
		return
	}
	if exists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Server error"})
		return
	}

	_, err = db.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
		req.Username, req.Email, string(hash), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User registered successfully"})
}

func LoginHandler(c *gin.Context, db *sql.DB) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	var id int
	var hash string
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", req.Email).Scan(&id, &hash)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	// Generate JWT
	tokenString, err := generateJWT(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful", "token": tokenString})
}

// JWT Middleware to protect routes
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization header required"})
			return
		}

		// Remove Bearer prefix if present
		parts := strings.Split(tokenString, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		} else if len(parts) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid Authorization header format"})
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
			return
		}

		// Store user ID in context
		c.Set("user_id", claims.UserID)

		c.Next() // proceed to the next handler
	}
}
