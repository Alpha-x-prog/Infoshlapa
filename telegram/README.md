# Telegram News Bot

## Установка и настройка

1. Установите Python зависимости:
```bash
cd scripts
pip install -r requirements.txt
```

2. Настройте конфигурацию:
```bash
# Скопируйте пример конфигурации
cp config.py.example config.py

# Отредактируйте config.py и укажите свои API ключи
# API_ID и API_HASH можно получить на https://my.telegram.org
```

3. Первый запуск Python скриптов:
```bash
# Добавление канала
python scripts/add_channel.py "https://t.me/your_channel"

# Тестовый запуск получения сообщений
python scripts/fetch_messages.py
```

## Запуск через Go

1. В основном приложении импортируйте пакет telegram:
```go
import "your_project/telegram"
```

2. Инициализируйте фоновый процесс получения сообщений:
```go
func main() {
    // Запуск фонового процесса получения сообщений
    telegram.StartMessageFetcher()
    
    // ... остальной код ...
}
```

3. Для добавления нового канала:
```go
err := telegram.AddChannel("https://t.me/your_channel")
if err != nil {
    // обработка ошибки
}
```

## Структура проекта

```
telegram/
├── scripts/
│   ├── add_channel.py      # Скрипт добавления канала
│   ├── fetch_messages.py   # Скрипт получения сообщений
│   ├── config.py          # Конфигурация (не в репозитории)
│   ├── config.py.example  # Пример конфигурации
│   └── requirements.txt   # Python зависимости
├── channel_adder.go       # Go обертка для добавления каналов
└── message_fetcher.go     # Go обертка для получения сообщений
```

## Важные замечания

1. При первом запуске Python скриптов потребуется авторизация в Telegram
2. Файл сессии (`*.session`) создастся автоматически
3. Логи сохраняются в `telegram_bot.log`
4. Убедитесь, что Python установлен в системе и доступен в PATH

## Функциональность

- Автоматическое создание необходимых таблиц в базе данных
- Мониторинг сообщений из указанных каналов
- Сохранение сообщений в базу данных
- Логирование всех действий в файл `telegram_bot.log`

## Структура базы данных

### Таблица telegram_channels
- channel_id: ID канала
- channel_username: Имя пользователя канала
- channel_title: Название канала
- last_message_id: ID последнего обработанного сообщения
- added_at: Дата добавления канала
- is_active: Статус активности канала

### Таблица telegram_messages
- message_id: ID сообщения
- channel_id: ID канала
- message_text: Текст сообщения
- message_date: Дата сообщения
- media_url: URL медиафайла (если есть)

## Безопасность

- Все учетные данные хранятся в файле `.env`
- Файл `.env` добавлен в `.gitignore` и не должен попадать в систему контроля версий
- Рекомендуется использовать разные учетные данные для разработки и продакшена 