import sys
import logging
from telethon import TelegramClient, events
from config import API_ID, API_HASH, SESSION_NAME, LOG_LEVEL, LOG_FILE, PHONE, PASSWORD_2FA
from database import Database
import os
import aiofiles
from datetime import datetime

# Настройка логирования
logging.basicConfig(
    level=getattr(logging, LOG_LEVEL),
    filename=LOG_FILE,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
# Добавляем вывод в консоль
console = logging.StreamHandler()
console.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
console.setFormatter(formatter)
logging.getLogger('').addHandler(console)

logger = logging.getLogger(__name__)

async def get_media_url(message, client):
    """Получение и сохранение медиафайла из сообщения"""
    try:
        if message.media:
            logger.info(f"Found media in message {message.id}")
            # Создаем имя файла на основе ID сообщения и текущей даты
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            filename = f"{message.id}_{timestamp}"
            
            if hasattr(message.media, 'photo'):
                logger.info(f"Processing photo in message {message.id}")
                # Фотографии
                file_path = f"telegram/media/{filename}.jpg"
                logger.info(f"Downloading photo to {file_path}")
                await client.download_media(message, file_path)
                logger.info(f"Successfully downloaded photo to {file_path}")
                return file_path
            elif hasattr(message.media, 'document'):
                # Проверяем, является ли документ изображением
                if message.media.document.mime_type and message.media.document.mime_type.startswith('image/'):
                    logger.info(f"Processing image document in message {message.id}")
                    # Определяем расширение файла из MIME-типа
                    ext = message.media.document.mime_type.split('/')[-1]
                    file_path = f"telegram/media/{filename}.{ext}"
                    logger.info(f"Downloading document to {file_path}")
                    await client.download_media(message, file_path)
                    logger.info(f"Successfully downloaded document to {file_path}")
                    return file_path
            elif hasattr(message.media, 'webpage'):
                # Веб-страницы с фото
                if message.media.webpage.photo:
                    logger.info(f"Processing webpage photo in message {message.id}")
                    file_path = f"telegram/media/{filename}.jpg"
                    logger.info(f"Downloading webpage photo to {file_path}")
                    await client.download_media(message, file_path)
                    logger.info(f"Successfully downloaded webpage photo to {file_path}")
                    return file_path
            else:
                logger.info(f"Unknown media type in message {message.id}")
    except Exception as e:
        logger.error(f"Error getting media URL for message {message.id}: {str(e)}", exc_info=True)
    return None

async def add_channel(channel_url):
    try:
        logger.info(f"Starting to add channel: {channel_url}")
        
        # Создаем клиент
        logger.info("Creating Telegram client...")
        client = TelegramClient(SESSION_NAME, API_ID, API_HASH)
        
        # Запускаем клиент с авторизацией
        logger.info("Starting client with authentication...")
        await client.start(phone=PHONE, password=PASSWORD_2FA)
        logger.info("Client started successfully")
        
        # Получаем информацию о канале
        logger.info(f"Getting channel info for: {channel_url}")
        channel = await client.get_entity(channel_url)
        logger.info(f"Successfully found channel: {channel.title} (ID: {channel.id})")
        
        # Получаем последние 100 сообщений
        logger.info("Fetching last 100 messages...")
        messages = await client.get_messages(channel, limit=100)
        logger.info(f"Found {len(messages)} messages")
        
        # Сохраняем канал и сообщения в базу данных
        logger.info("Connecting to database...")
        db = Database()
        
        # Сохраняем канал
        if db.add_channel(channel.id, channel.username, channel.title):
            logger.info(f"Channel {channel.title} saved to database")
            
            # Сохраняем сообщения
            last_message_id = 0
            for message in messages:
                if message.id > last_message_id:
                    last_message_id = message.id
                
                # Получаем URL медиа
                media_url = await get_media_url(message, client)
                
                # Сохраняем сообщение
                db.add_message(
                    message_id=message.id,
                    channel_id=channel.id,
                    text=message.text,
                    date=message.date,
                    media_url=media_url
                )
            
            # Обновляем last_message_id в таблице channels
            db.update_last_message_id(channel.id, last_message_id)
            logger.info(f"Saved {len(messages)} messages and updated last_message_id to {last_message_id}")
        else:
            logger.error(f"Failed to save channel {channel.title} to database")
        
        await client.disconnect()
        logger.info("Client disconnected")
        return True
    except Exception as e:
        logger.error(f"Error adding channel: {str(e)}", exc_info=True)
        return False

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python add_channel.py <channel_url>")
        sys.exit(1)
    
    channel_url = sys.argv[1]
    logger.info(f"Script started with channel URL: {channel_url}")
    
    import asyncio
    try:
        asyncio.run(add_channel(channel_url))
    except Exception as e:
        logger.error(f"Script failed with error: {str(e)}", exc_info=True) 