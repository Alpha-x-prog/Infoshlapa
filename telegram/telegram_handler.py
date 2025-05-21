import os
import sqlite3
from telethon import TelegramClient, events
from telethon.tl.types import Channel
import asyncio
import logging
from datetime import datetime
import json
from dotenv import load_dotenv

# Загрузка переменных окружения из .env файла
load_dotenv()

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    filename='telegram_bot.log'
)
logger = logging.getLogger(__name__)

# Путь к базе данных
DB_PATH = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'news.db'))

class TelegramHandler:
    def __init__(self):
        self.client = None
        self.db = None
        self.api_id = os.getenv('TELEGRAM_API_ID')
        self.api_hash = os.getenv('TELEGRAM_API_HASH')
        self.phone = os.getenv('TELEGRAM_PHONE')
        self.password = os.getenv('TELEGRAM_PASSWORD')

    async def init_client(self):
        """Инициализация клиента Telegram"""
        if not all([self.api_id, self.api_hash, self.phone, self.password]):
            raise ValueError("Необходимо установить все переменные окружения: TELEGRAM_API_ID, TELEGRAM_API_HASH, TELEGRAM_PHONE, TELEGRAM_PASSWORD")
        
        self.client = TelegramClient('news_bot_session', self.api_id, self.api_hash)
        await self.client.start(phone=self.phone, password=self.password)
        logger.info("Telegram клиент инициализирован")

    def init_db(self):
        """Инициализация подключения к базе данных"""
        self.db = sqlite3.connect(DB_PATH)
        self.create_tables()
        logger.info("Подключение к базе данных установлено")

    def create_tables(self):
        """Создание необходимых таблиц в базе данных"""
        cursor = self.db.cursor()
        
        # Таблица для хранения информации о каналах
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS telegram_channels (
                channel_id INTEGER PRIMARY KEY AUTOINCREMENT,
                channel_username TEXT UNIQUE,
                channel_title TEXT,
                last_message_id INTEGER DEFAULT 0,
                added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                is_active BOOLEAN DEFAULT 1
            )
        ''')

        # Таблица для хранения сообщений из каналов
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS telegram_messages (
                message_id INTEGER PRIMARY KEY,
                channel_id INTEGER,
                message_text TEXT,
                message_date TIMESTAMP,
                media_url TEXT,
                FOREIGN KEY (channel_id) REFERENCES telegram_channels(channel_id)
            )
        ''')
        
        self.db.commit()

    def add_channel(self, username):
        """Добавление нового канала в базу данных"""
        try:
            cursor = self.db.cursor()
            cursor.execute(
                "INSERT OR IGNORE INTO telegram_channels (channel_username) VALUES (?)",
                (username,)
            )
            self.db.commit()
            logger.info(f"Канал {username} добавлен в базу данных")
            return True
        except Exception as e:
            logger.error(f"Ошибка при добавлении канала {username}: {e}")
            return False

    def get_active_channels(self):
        """Получение списка активных каналов"""
        try:
            cursor = self.db.cursor()
            cursor.execute("SELECT channel_username FROM telegram_channels WHERE is_active = 1")
            return [row[0] for row in cursor.fetchall()]
        except Exception as e:
            logger.error(f"Ошибка при получении списка каналов: {e}")
            return []

    async def get_channel_messages(self, channel_username, limit=100):
        """Получение сообщений из канала"""
        try:
            channel = await self.client.get_entity(channel_username)
            messages = await self.client.get_messages(channel, limit=limit)
            return messages
        except Exception as e:
            logger.error(f"Ошибка при получении сообщений из канала {channel_username}: {e}")
            return []

    def save_message(self, channel_id, message):
        """Сохранение сообщения в базу данных"""
        try:
            cursor = self.db.cursor()
            media_url = None
            if message.media:
                # Здесь можно добавить логику сохранения медиафайлов
                pass

            cursor.execute('''
                INSERT OR IGNORE INTO telegram_messages 
                (message_id, channel_id, message_text, message_date, media_url)
                VALUES (?, ?, ?, ?, ?)
            ''', (
                message.id,
                channel_id,
                message.text,
                message.date,
                media_url
            ))
            self.db.commit()
            logger.info(f"Сообщение {message.id} сохранено в базу данных")
        except Exception as e:
            logger.error(f"Ошибка при сохранении сообщения: {e}")

    async def process_channel(self, channel_username):
        """Обработка канала и сохранение новых сообщений"""
        try:
            cursor = self.db.cursor()
            # Получаем ID канала и последнее обработанное сообщение
            cursor.execute(
                "SELECT channel_id, last_message_id FROM telegram_channels WHERE channel_username = ?",
                (channel_username,)
            )
            channel_data = cursor.fetchone()
            if not channel_data:
                return
            
            channel_id, last_message_id = channel_data
            
            # Получаем новые сообщения
            messages = await self.get_channel_messages(channel_username)
            new_messages = [msg for msg in messages if msg.id > last_message_id]

            if new_messages:
                # Сохраняем новые сообщения
                for message in reversed(new_messages):
                    self.save_message(channel_id, message)

                # Обновляем ID последнего сообщения
                cursor.execute(
                    "UPDATE telegram_channels SET last_message_id = ? WHERE channel_id = ?",
                    (new_messages[0].id, channel_id)
                )
                self.db.commit()
                logger.info(f"Обработано {len(new_messages)} новых сообщений из канала {channel_username}")

        except Exception as e:
            logger.error(f"Ошибка при обработке канала {channel_username}: {e}")

    async def monitor_channels(self):
        """Мониторинг всех активных каналов"""
        while True:
            channels = self.get_active_channels()
            for channel in channels:
                await self.process_channel(channel)
            await asyncio.sleep(1800)  # Проверка каждые 30 минут

    async def run(self):
        """Основной метод запуска"""
        await self.init_client()
        self.init_db()
        
        try:
            await self.monitor_channels()
        except KeyboardInterrupt:
            logger.info("Остановка Telegram клиента...")
        finally:
            await self.client.disconnect()
            self.db.close()

if __name__ == "__main__":
    handler = TelegramHandler()
    asyncio.run(handler.run()) 