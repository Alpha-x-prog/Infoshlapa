import sys
import logging
from telethon import TelegramClient, events
from config import API_ID, API_HASH, SESSION_NAME, LOG_LEVEL, LOG_FILE, PHONE, PASSWORD_2FA
from database import Database
import os
import aiofiles
from datetime import datetime
import requests
from dotenv import load_dotenv
import urllib3

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
        
        prompt = f"Сделай краткое изложение следующего текста в 3-4 предложения: {text}"
        
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
        
        # Получаем последние 5 сообщений
        logger.info("Fetching last 5 messages...")
        messages = await client.get_messages(channel, limit=20)  # Увеличиваем лимит для пропуска пустых сообщений
        logger.info(f"Found {len(messages)} messages")
        
        # Получаем фото канала для использования по умолчанию
        channel_photo = None
        try:
            if channel.photo:
                logger.info("Downloading channel photo...")
                timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
                channel_photo = f"telegram/media/channel_{channel.id}_{timestamp}.jpg"
                await client.download_media(channel.photo, channel_photo)
                logger.info(f"Channel photo downloaded: {channel_photo}")
        except Exception as e:
            logger.error(f"Error downloading channel photo: {str(e)}")
        
        # Сохраняем канал и сообщения в базу данных
        logger.info("Connecting to database...")
        db = Database()
        
        # Сохраняем канал
        if db.add_channel(channel.id, channel.username, channel.title):
            logger.info(f"Channel {channel.title} saved to database")
            
            # Сохраняем сообщения
            last_message_id = 0
            saved_messages_count = 0
            max_messages = 5  # Желаемое количество сохраненных сообщений
            
            for message in messages:
                if message.id > last_message_id:
                    last_message_id = message.id
                
                logger.info(f"Processing message {message.id}")
                
                # Пропускаем сообщения без текста
                if not message.text or message.text.strip() == "":
                    logger.info(f"Skipping message {message.id} - no text content")
                    continue
                
                logger.info(f"Message {message.id} text: {message.text[:100]}...")  # Логируем первые 100 символов текста
                
                # Получаем URL медиа
                media_url = await get_media_url(message, client)
                
                # Если нет медиа в сообщении, используем фото канала
                if not media_url:
                    if channel_photo:
                        logger.info(f"Using channel photo as default media for message {message.id}")
                        media_url = channel_photo
                    else:
                        logger.info(f"No media found for message {message.id} and no channel photo available")
                
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
                    channel_id=channel.id,
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
            
            # Обновляем last_message_id в таблице channels
            db.update_last_message_id(channel.id, last_message_id)
            logger.info(f"Saved {saved_messages_count} messages out of {len(messages)} processed and updated last_message_id to {last_message_id}")
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