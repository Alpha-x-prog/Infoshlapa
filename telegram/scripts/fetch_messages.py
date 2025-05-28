import logging
from telethon import TelegramClient
import json
import os
from config import (
    API_ID, API_HASH, SESSION_NAME, 
    MAX_MESSAGES_PER_CHANNEL, LOG_LEVEL, LOG_FILE,
    PHONE, PASSWORD_2FA
)
from database import Database
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

async def fetch_messages():
    try:
        # Создаем клиент
        client = TelegramClient(SESSION_NAME, API_ID, API_HASH)
        
        # Запускаем клиент с авторизацией
        await client.start(phone=PHONE, password=PASSWORD_2FA)
        
        # Получаем список каналов из базы данных
        db = Database()
        channels = db.get_channels()
        logger.info(f"Found {len(channels)} active channels")
        
        for channel_id, username, title, last_message_id in channels:
            try:
                logger.info(f"Processing channel: {title} (last_message_id: {last_message_id})")
                
                # Получаем канал
                channel = await client.get_entity(username)
                
                # Получаем новые сообщения
                messages = await client.get_messages(
                    channel,
                    min_id=last_message_id,  # Получаем только сообщения новее last_message_id
                    limit=10
                )
                
                if messages:
                    logger.info(f"Found {len(messages)} new messages in {title}")
                    
                    # Сохраняем новые сообщения
                    new_last_message_id = last_message_id
                    for message in messages:
                        if message.id > new_last_message_id:
                            new_last_message_id = message.id
                        
                        # Получаем URL медиа
                        media_url = await get_media_url(message, client)
                        
                        # Сохраняем сообщение
                        db.add_message(
                            message_id=message.id,
                            channel_id=channel_id,
                            text=message.text,
                            date=message.date,
                            media_url=media_url
                        )
                    
                    # Обновляем last_message_id
                    if new_last_message_id > last_message_id:
                        db.update_last_message_id(channel_id, new_last_message_id)
                        logger.info(f"Updated last_message_id to {new_last_message_id} for {title}")
                else:
                    logger.info(f"No new messages in {title}")
                
            except Exception as e:
                logger.error(f"Error processing channel {title}: {str(e)}")
        
        await client.disconnect()
        return True
    except Exception as e:
        logger.error(f"Error in fetch_messages: {str(e)}")
        return False

if __name__ == "__main__":
    import asyncio
    asyncio.run(fetch_messages()) 