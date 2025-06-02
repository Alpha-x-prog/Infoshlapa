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
import requests
from dotenv import load_dotenv
import urllib3
import sys

# Загружаем переменные окружения из .env файла
load_dotenv()

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

# Проверяем наличие необходимых переменных окружения
required_env_vars = {
    'GEMINI_API_KEY': 'API ключ для Gemini (https://makersuite.google.com/app/apikey)',
    'TELEGRAM_API_ID': 'Telegram API ID (https://my.telegram.org)',
    'TELEGRAM_API_HASH': 'Telegram API Hash (https://my.telegram.org)',
    'PROXY_URL': 'URL прокси-сервера',
    'PROXY_USER': 'Имя пользователя прокси',
    'PROXY_PASSWORD': 'Пароль прокси'
}

missing_vars = []
for var, description in required_env_vars.items():
    if not os.getenv(var):
        missing_vars.append(f"{var} - {description}")

if missing_vars:
    error_msg = "Отсутствуют необходимые переменные окружения:\n" + "\n".join(missing_vars)
    logger.error(error_msg)
    print(error_msg)
    sys.exit(1)

def get_gemini_summary(text):
    """Получает краткое изложение текста от Gemini API"""
    try:
        api_key = os.getenv('GEMINI_API_KEY')
        if not api_key:
            error_msg = "GEMINI_API_KEY не найден в переменных окружения. Получите ключ на https://makersuite.google.com/app/apikey"
            logger.error(error_msg)
            print(error_msg)
            return None

        # Настройка прокси
        proxy_url = os.getenv('PROXY_URL')
        proxy_user = os.getenv('PROXY_USER')
        proxy_pass = os.getenv('PROXY_PASSWORD')
        
        proxies = {
            'http': f'http://{proxy_user}:{proxy_pass}@{proxy_url}',
            'https': f'http://{proxy_user}:{proxy_pass}@{proxy_url}'
        }
        
        # Отключаем предупреждения о небезопасных SSL-сертификатах
        urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
        
        url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key={api_key}"
        
        prompt = f"Сделай краткое изложение следующего текста в 2-3 предложения, сохраняя основной смысл:\n\n{text}"
        
        payload = {
            "contents": [{
                "parts": [{
                    "text": prompt
                }]
            }]
        }
        
        headers = {
            "Content-Type": "application/json"
        }
        
        # Используем прокси и отключаем проверку SSL
        response = requests.post(
            url, 
            json=payload, 
            headers=headers,
            proxies=proxies,
            verify=False
        )
        response.raise_for_status()
        
        result = response.json()
        if 'candidates' in result and len(result['candidates']) > 0:
            return result['candidates'][0]['content']['parts'][0]['text']
        return None
    except Exception as e:
        logger.error(f"Error getting Gemini summary: {str(e)}")
        return None

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
                    limit=20  # Увеличиваем лимит, чтобы иметь запас для пропуска сообщений без текста
                )
                
                if messages:
                    logger.info(f"Found {len(messages)} new messages in {title}")
                    
                    # Сохраняем новые сообщения
                    new_last_message_id = last_message_id
                    saved_messages_count = 0
                    max_messages = 5  # Желаемое количество сохраненных сообщений
                    
                    for message in messages:
                        if message.id > new_last_message_id:
                            new_last_message_id = message.id
                        
                        # Пропускаем сообщения без текста
                        if not message.text:
                            logger.info(f"Skipping message {message.id} - no text content")
                            continue
                        
                        logger.info(f"Processing message {message.id}")
                        logger.info(f"Message {message.id} text: {message.text[:100]}...")  # Логируем первые 100 символов текста
                        
                        # Получаем URL медиа
                        media_url = await get_media_url(message, client)
                        if media_url:
                            logger.info(f"Media URL for message {message.id}: {media_url}")
                        
                        # Генерируем краткое изложение
                        logger.info(f"Generating summary for message {message.id}")
                        summary = get_gemini_summary(message.text)
                        if summary:
                            logger.info(f"Generated summary for message {message.id}: {summary[:100]}...")  # Логируем первые 100 символов summary
                        else:
                            logger.error(f"Failed to generate summary for message {message.id}")
                        
                        # Сохраняем сообщение
                        logger.info(f"Saving message {message.id} to database")
                        success = db.add_message(
                            message_id=message.id,
                            channel_id=channel_id,
                            text=message.text,
                            date=message.date,
                            media_url=media_url,
                            summary=summary
                        )
                        if success:
                            logger.info(f"Successfully saved message {message.id}")
                            saved_messages_count += 1
                        else:
                            logger.error(f"Failed to save message {message.id}")
                        
                        # Проверяем, достигли ли мы нужного количества сообщений
                        if saved_messages_count >= max_messages:
                            logger.info(f"Reached target of {max_messages} saved messages")
                            break
                    
                    # Обновляем last_message_id
                    if new_last_message_id > last_message_id:
                        db.update_last_message_id(channel_id, new_last_message_id)
                        logger.info(f"Updated last_message_id to {new_last_message_id} for {title}")
                        logger.info(f"Saved {saved_messages_count} messages out of {len(messages)} processed")
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