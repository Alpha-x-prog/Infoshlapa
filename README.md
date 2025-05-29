# News API Documentation

## Аутентификация

### Регистрация нового пользователя
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "your_password"
  }'
```

Ответ:
```json
{
    "success": true,
    "token": "YOUR_JWT_TOKEN",
    "user": {
        "id": 1,
        "email": "test@example.com"
    },
    "message": "Пользователь успешно зарегистрирован"
}
```

### Вход пользователя
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "your_password"
  }'
```

Ответ:
```json
{
    "success": true,
    "token": "YOUR_JWT_TOKEN",
    "user": {
        "id": 1,
        "email": "test@example.com"
    },
    "message": "Успешная авторизация"
}
```

## Новости

### Получение списка новостей
```bash
# Базовый запрос
curl -X GET http://localhost:8080/api/news

# С параметрами
curl -X GET "http://localhost:8080/api/news?category=technology&offset=0"
```

Ответ:
```json
[
    {
        "article_id": "123",
        "title": "Example News",
        "link": "https://example.com/news",
        "keywords": ["tech", "news"],
        "creator": ["John Doe"],
        "video_url": "",
        "description": "News description",
        "content": "Full content",
        "publishedAt": "2024-03-21T10:00:00Z",
        "urlToImage": "https://example.com/image.jpg",
        "source_id": "example",
        "source_name": "Example News",
        "url": "https://example.com",
        "language": "en",
        "country": "US",
        "tags": "technology",
        "sentiment": "positive"
    }
]
```

## Закладки

### Добавление закладки
```bash
curl -X POST http://localhost:8080/api/protected/bookmarks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "newsId": "123"
  }'
```

Ответ:
```json
{
    "message": "Bookmark added successfully"
}
```

### Получение закладок
```bash
curl -X GET http://localhost:8080/api/protected/bookmarks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Ответ:
```json
[
    {
        "article_id": "123",
        "title": "Example News",
        "link": "https://example.com/news",
        "keywords": ["tech", "news"],
        "creator": ["John Doe"],
        "video_url": "",
        "description": "News description",
        "content": "Full content",
        "publishedAt": "2024-03-21T10:00:00Z",
        "urlToImage": "https://example.com/image.jpg",
        "source_id": "example",
        "source_name": "Example News",
        "url": "https://example.com",
        "language": "en",
        "country": "US",
        "tags": "technology",
        "sentiment": "positive"
    }
]
```

### Удаление закладки
```bash
curl -X DELETE http://localhost:8080/api/protected/bookmarks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "newsId": "123"
  }'
```

Ответ:
```json
{
    "message": "Bookmark removed successfully"
}
```

## Telegram Каналы

### Добавление канала
```bash
curl -X POST http://localhost:8080/api/protected/channels \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "channel_url": "https://t.me/example_channel",
    "channel_name": "Example Channel"
  }'
```

Ответ:
```json
{
    "success": true,
    "message": "Channel added successfully"
}
```

### Получение каналов
```bash
curl -X GET http://localhost:8080/api/protected/channels \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Ответ:
```json
{
    "success": true,
    "channels": [
        {
            "url": "https://t.me/example_channel",
            "name": "Example Channel",
            "created_at": "2024-03-21T10:30:00Z"
        }
    ]
}
```

### Удаление канала
```bash
curl -X DELETE "http://localhost:8080/api/protected/channels?channel_url=https://t.me/example_channel" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Ответ:
```json
{
    "success": true,
    "message": "Channel removed successfully"
}
```

## AI Запросы

### Запрос к Gemini AI
```bash
curl -X POST http://localhost:8080/api/ask \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "What is the latest news about technology?"
  }'
```

Ответ:
```json
{
    "content": "AI response here..."
}
```

## Административные функции

### Удаление всех пользователей
```bash
curl -X DELETE http://localhost:8080/api/protected/users/all \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Ответ:
```json
{
    "success": true,
    "message": "All users have been deleted"
}
```

## Коды ошибок

- 200: Успешный запрос
- 400: Неверный формат данных
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Ресурс не найден
- 500: Внутренняя ошибка сервера

## Примеры ошибок

### Неверный формат данных
```json
{
    "error": "Неверный формат данных"
}
```

### Пользователь не авторизован
```json
{
    "error": "Unauthorized"
}
```

### Пользователь уже существует
```json
{
    "error": "Пользователь с таким email уже существует"
}
```

### Неверный пароль
```json
{
    "error": "Неверный пароль"
}
```

## Примечания

1. Все защищенные маршруты требуют JWT токен в заголовке Authorization
2. Токен получается при регистрации или входе
3. Токен действителен 7 дней
4. Для всех запросов используется формат JSON
5. Все даты возвращаются в формате ISO 8601
