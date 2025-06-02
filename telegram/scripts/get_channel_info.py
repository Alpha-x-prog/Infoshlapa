import sys
import logging
from telethon import TelegramClient
from telethon.tl.functions.channels import GetFullChannelRequest
from telethon.tl.types import Channel, Chat
import os
from dotenv import load_dotenv
import io

# Настройка кодировки для stdout
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Загрузка переменных окружения
load_dotenv()

# Конфигурация
API_ID = os.getenv('TELEGRAM_API_ID')
API_HASH = os.getenv('TELEGRAM_API_HASH')
PHONE = os.getenv('TELEGRAM_PHONE')
PASSWORD_2FA = os.getenv('TELEGRAM_PASSWORD_2FA')
SESSION_NAME = 'news_bot_session'

async def get_channel_info(channel_url):
    try:
        logger.info(f"Starting to get info for channel: {channel_url}")
        
        # Создаем клиент
        logger.info("Creating Telegram client...")
        client = TelegramClient(SESSION_NAME, API_ID, API_HASH)
        
        # Запускаем клиент с авторизацией
        logger.info("Starting client with authentication...")
        await client.start(phone=PHONE, password=PASSWORD_2FA)
        logger.info("Client started successfully")
        
        # Получаем информацию о канале
        logger.info(f"Getting channel info for: {channel_url}")
        entity = await client.get_entity(channel_url)
        
        channel_title = ""
        
        # Пробуем получить название канала разными способами
        if isinstance(entity, Channel):
            try:
                # Пробуем получить полную информацию о канале
                full_channel = await client(GetFullChannelRequest(entity))
                channel_title = full_channel.chats[0].title
            except Exception as e:
                logger.warning(f"Failed to get full channel info: {str(e)}")
                # Если не удалось получить полную информацию, используем базовое название
                channel_title = entity.title
        elif isinstance(entity, Chat):
            channel_title = entity.title
        else:
            # Если сущность не является ни каналом, ни чатом, используем username как fallback
            channel_title = channel_url.split('/')[-1]
        
        # Если название все еще пустое, используем username
        if not channel_title:
            channel_title = channel_url.split('/')[-1]
        
        logger.info(f"Successfully got channel title: {channel_title}")
        print(channel_title.encode('utf-8').decode('utf-8'))  # Явно обрабатываем кодировку
        
        await client.disconnect()
        logger.info("Client disconnected")
        return True
    except Exception as e:
        logger.error(f"Error getting channel info: {str(e)}", exc_info=True)
        # В случае ошибки используем username как fallback
        fallback_title = channel_url.split('/')[-1]
        logger.info(f"Using fallback title: {fallback_title}")
        print(fallback_title.encode('utf-8').decode('utf-8'))  # Явно обрабатываем кодировку
        return False

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python get_channel_info.py <channel_url>")
        sys.exit(1)
    
    channel_url = sys.argv[1]
    logger.info(f"Script started with channel URL: {channel_url}")
    
    import asyncio
    try:
        asyncio.run(get_channel_info(channel_url))
    except Exception as e:
        logger.error(f"Script failed with error: {str(e)}", exc_info=True)
        # В случае ошибки используем username как fallback
        fallback_title = channel_url.split('/')[-1]
        print(fallback_title.encode('utf-8').decode('utf-8'))  # Явно обрабатываем кодировку 