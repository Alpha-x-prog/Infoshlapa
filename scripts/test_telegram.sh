#!/bin/bash

# Проверяем наличие Python и pip
if ! command -v python3 &> /dev/null; then
    echo "Python3 не установлен"
    exit 1
fi

if ! command -v pip3 &> /dev/null; then
    echo "pip3 не установлен"
    exit 1
fi

# Устанавливаем зависимости Python
echo "Установка зависимостей Python..."
pip3 install -r telegram/requirements.txt

# Проверяем наличие .env файла
if [ ! -f telegram/.env ]; then
    echo "Создание .env файла..."
    echo "TELEGRAM_API_ID=your_api_id" > telegram/.env
    echo "TELEGRAM_API_HASH=your_api_hash" >> telegram/.env
    echo "TELEGRAM_PHONE=your_phone_number" >> telegram/.env
    echo "TELEGRAM_PASSWORD=your_2fa_password" >> telegram/.env
    echo "Пожалуйста, отредактируйте telegram/.env и добавьте ваши учетные данные Telegram"
    exit 1
fi

# Запускаем тесты
echo "Запуск тестов..."
go test ./telegram -v

# Запускаем тестовый скрипт
echo "Запуск тестового скрипта..."
go run cmd/test_telegram/main.go

# Очистка
echo "Очистка тестовых файлов..."
rm -f test_news.db test_telegram_bot.log 