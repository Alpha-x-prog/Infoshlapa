import sqlite3
import logging
from datetime import datetime
import os

logger = logging.getLogger(__name__)

class Database:
    def __init__(self, db_file='news.db'):
        # Получаем путь к корневой директории проекта
        current_dir = os.path.dirname(os.path.abspath(__file__))  # текущая директория скрипта
        project_root = os.path.abspath(os.path.join(current_dir, '..', '..'))  # поднимаемся на два уровня вверх
        self.db_file = os.path.join(project_root, db_file)  # полный путь к базе данных
        logger.info(f"Using database file: {self.db_file}")
        self.init_db()

    def init_db(self):
        """Инициализация базы данных и создание таблиц"""
        try:
            conn = sqlite3.connect(self.db_file)
            cursor = conn.cursor()

            # Создаем таблицу для каналов, если она не существует
            cursor.execute('''
                CREATE TABLE IF NOT EXISTS telegram_channels (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    channel_id INTEGER UNIQUE,
                    channel_username TEXT,
                    channel_title TEXT,
                    last_message_id INTEGER DEFAULT 0,
                    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    is_active BOOLEAN DEFAULT 1
                )
            ''')

            # Создаем таблицу для сообщений, если она не существует
            cursor.execute('''
                CREATE TABLE IF NOT EXISTS telegram_messages (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    message_id INTEGER,
                    channel_id INTEGER,
                    message_text TEXT,
                    message_date TIMESTAMP,
                    media_url TEXT,
                    FOREIGN KEY (channel_id) REFERENCES telegram_channels (channel_id)
                )
            ''')

            conn.commit()
            logger.info("Database initialized successfully")
        except Exception as e:
            logger.error(f"Error initializing database: {str(e)}")
            raise
        finally:
            conn.close()

    def add_channel(self, channel_id, username, title):
        """Добавление нового канала"""
        try:
            conn = sqlite3.connect(self.db_file)
            cursor = conn.cursor()
            
            cursor.execute('''
                INSERT OR REPLACE INTO telegram_channels 
                (channel_id, channel_username, channel_title, last_message_id, added_at, is_active)
                VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP, 1)
            ''', (channel_id, username, title))
            
            conn.commit()
            logger.info(f"Channel {title} added successfully")
            return True
        except Exception as e:
            logger.error(f"Error adding channel: {str(e)}")
            return False
        finally:
            conn.close()

    def update_last_message_id(self, channel_id, last_message_id):
        """Обновление last_message_id для канала"""
        try:
            conn = sqlite3.connect(self.db_file)
            cursor = conn.cursor()
            
            cursor.execute('''
                UPDATE telegram_channels
                SET last_message_id = ?
                WHERE channel_id = ?
            ''', (last_message_id, channel_id))
            
            conn.commit()
            logger.info(f"Updated last_message_id to {last_message_id} for channel {channel_id}")
            return True
        except Exception as e:
            logger.error(f"Error updating last_message_id: {str(e)}")
            return False
        finally:
            conn.close()

    def get_channels(self):
        """Получение списка активных каналов"""
        try:
            conn = sqlite3.connect(self.db_file)
            cursor = conn.cursor()
            
            cursor.execute('''
                SELECT channel_id, channel_username, channel_title, last_message_id
                FROM telegram_channels
                WHERE is_active = 1
            ''')
            
            return cursor.fetchall()
        except Exception as e:
            logger.error(f"Error getting channels: {str(e)}")
            return []
        finally:
            conn.close()

    def add_message(self, message_id, channel_id, text, date, media_url=None):
        """Добавление нового сообщения"""
        try:
            conn = sqlite3.connect(self.db_file)
            cursor = conn.cursor()
            
            cursor.execute('''
                INSERT INTO telegram_messages 
                (message_id, channel_id, message_text, message_date, media_url)
                VALUES (?, ?, ?, ?, ?)
            ''', (message_id, channel_id, text, date, media_url))
            
            conn.commit()
            return True
        except Exception as e:
            logger.error(f"Error adding message: {str(e)}")
            return False
        finally:
            conn.close() 