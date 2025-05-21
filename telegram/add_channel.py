from telegram_handler import TelegramHandler
import asyncio

async def add_channel(username):
    handler = TelegramHandler()
    handler.init_db()
    handler.add_channel(username)
    print(f"Канал {username} добавлен для мониторинга")

if __name__ == "__main__":
    # Пример добавления канала
    channel_username = input("Введите имя пользователя канала (например, @channel_name): ")
    if channel_username.startswith('@'):
        channel_username = channel_username[1:]
    asyncio.run(add_channel(channel_username)) 