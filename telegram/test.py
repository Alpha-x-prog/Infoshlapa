from telegram_handler import TelegramHandler
import asyncio
import sqlite3

async def test_connection():
    print("Тестирование подключения к Telegram...")
    handler = TelegramHandler()
    await handler.init_client()
    print("Подключение успешно установлено!")

def test_database():
    print("\nТестирование базы данных...")
    conn = sqlite3.connect('news.db')
    cursor = conn.cursor()
    
    # Проверяем таблицу каналов
    cursor.execute("SELECT * FROM telegram_channels")
    channels = cursor.fetchall()
    print(f"Найдено каналов: {len(channels)}")
    for channel in channels:
        print(f"- {channel[1]} (ID: {channel[0]})")
    
    # Проверяем таблицу сообщений
    cursor.execute("SELECT COUNT(*) FROM telegram_messages")
    message_count = cursor.fetchone()[0]
    print(f"\nВсего сообщений в базе: {message_count}")
    
    conn.close()

if __name__ == "__main__":
    print("=== Тестирование Telegram интеграции ===")
    
    # Тестируем подключение
    asyncio.run(test_connection())
    
    # Тестируем базу данных
    test_database()
    
    print("\nТестирование завершено!") 