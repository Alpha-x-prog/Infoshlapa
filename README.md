# News API

Сервис для агрегации новостей с поддержкой Telegram-каналов и AI-анализа.

## Содержание

- [Требования](#требования)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [Запуск](#запуск)
- [API Документация](#api-документация)
- [Структура проекта](#структура-проекта)
- [Разработка](#разработка)

## Требования

- Go 1.24.0 или выше
- SQLite3
- Node.js 16+ (для фронтенда)
- Telegram Bot Token

## Установка

1. Клонируйте репозиторий:
```bash
git clone [URL репозитория]
cd newsAPI
```

2. Установите зависимости Go:
```bash
go mod download
```

3. Установите зависимости фронтенда:
```bash
cd shlapa
npm install
```

## Конфигурация

1. Создайте файл `.env` в корневой директории:
```env
# Database configuration
DB_PATH=news.db

# Server configuration
PORT=8080

# Telegram Bot configuration
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Paths (relative to project root)
STATIC_PATH=static
DIST_PATH=shlapa/dist
```

## Запуск

1. Запуск бэкенда:
```bash
go run main.go
```

2. Запуск фронтенда (в режиме разработки):
```bash
cd shlapa
npm run serve
```

3. Сборка фронтенда:
```bash
cd shlapa
npm run build
```

## API Документация

### Аутентификация

#### Регистрация
```bash
POST /api/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password"
}
```

#### Вход
```bash
POST /api/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password"
}
```

### Новости

#### Получение новостей
```bash
GET /api/news
GET /api/news?category=technology&offset=0
```

### Закладки

#### Добавление закладки
```bash
POST /api/protected/bookmarks
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
    "newsId": "123"
}
```

#### Получение закладок
```bash
GET /api/protected/bookmarks
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Удаление закладки
```bash
DELETE /api/protected/bookmarks
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
    "newsId": "123"
}
```

### Telegram Каналы

#### Добавление канала
```bash
POST /api/protected/channels
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
    "channel_url": "https://t.me/example",
    "channel_name": "Example Channel"
}
```

#### Получение каналов
```bash
GET /api/protected/channels
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Удаление канала
```bash
DELETE /api/protected/channels?channel_url=https://t.me/example
Authorization: Bearer YOUR_JWT_TOKEN
```

### AI Запросы

#### Запрос к Gemini AI
```bash
POST /api/ask
Content-Type: application/json

{
    "prompt": "Your question here"
}
```

## Структура проекта

```
newsAPI/
├── api/            # API endpoints
├── config/         # Конфигурация приложения
├── db/            # Работа с базой данных
├── handlers/      # Обработчики запросов
├── middleware/    # Промежуточное ПО
├── parser/        # Парсеры данных
├── telegram/      # Интеграция с Telegram
├── gemini/        # Интеграция с Gemini AI
├── collyan/       # Веб-скрапинг
├── static/        # Статические файлы
├── shlapa/        # Фронтенд (Vue.js)
├── image/         # Изображения
├── scripts/       # Скрипты
└── cmd/           # Исполняемые команды
```

## Разработка

### Коды ошибок

- 200: Успешный запрос
- 400: Неверный формат данных
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Ресурс не найден
- 500: Внутренняя ошибка сервера

### Логирование

Логи приложения записываются в:
- `telegram_bot.log` - логи Telegram бота
- Стандартный вывод для остальных логов

### Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...
```
