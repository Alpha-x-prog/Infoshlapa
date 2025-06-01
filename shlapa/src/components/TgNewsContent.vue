<template>
    <div class="tg-news-container">
        <div class="tg-news-header">
            <h1>Telegram новости</h1>

        </div>
        
        <div class="channel-form">
            <b-form @submit.prevent="addChannel" class="d-flex align-items-center">
                <b-form-input
                    v-model="channelUrl"
                    placeholder="Введите ссылку на канал (например, https://t.me/channel)"
                    class="mr-2 "
                    :state="channelUrlValidation"
                ></b-form-input>
                <b-button type="submit" variant="primary" :disabled="!channelUrlValidation">
                    Добавить канал
                </b-button>
            </b-form>

            <b-alert
                v-model="showAlert"
                :variant="alertVariant"
                dismissible
                class="mt-3"
            >
                {{ alertMessage }}
            </b-alert>
        </div>

        <div class="channels-list mt-4" v-if="channels.length > 0">
            <h3>Ваши каналы:</h3>
            <b-list-group>
                <b-list-group-item
                    v-for="channel in channels"
                    :key="channel.url"
                    class="d-flex justify-content-between align-items-center"
                >
                    <div>
                        <strong>{{ channel.name }}</strong>
                        <br>
                        <small class="text-muted">Добавлен: {{ formatDate(channel.created_at) }}</small>
                    </div>
                    <b-button variant="danger" size="sm" @click="removeChannel(channel.url)">
                        Удалить
                    </b-button>
                </b-list-group-item>
            </b-list-group>
        </div>
        <TimeBar/>
        <div class="messages-container mt-4">
            <div v-if="loading" class="text-center">
                <b-spinner label="Loading..."></b-spinner>
            </div>
            <div v-else-if="messages.length === 0" class="text-center">
                <p>Нет сообщений для отображения</p>
            </div>
            <div v-else class="messages-grid">
                <div v-for="message in messages" :key="message.message_id" class="message-card">
                    <b-card>
                        <b-card-title>{{ message.channel_title }}</b-card-title>
                        <b-card-sub-title class="mb-2">
                            {{ formatDate(message.date) }}
                        </b-card-sub-title>
                        <b-card-text>
                            {{ message.text }}
                        </b-card-text>
                        <!--<b-img v-if="message.media_url" :src="message.media_url" fluid alt="Message media"></b-img>-->
                        <template v-if="message.summary" #footer>
                            <small class="text-muted">{{ message.summary }}</small>
                        </template>
                    </b-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from '@/utils/axios';
import TimeBar from "@/components/TimeBar";

export default {
    name: 'TgNewsContent',
    data() {
        return {
            channelUrl: '',
            channels: [],
            messages: [],
            loading: false,
            showAlert: false,
            alertVariant: 'success',
            alertMessage: '',
            refreshInterval: null
        }
    },
    computed: {
        channelUrlValidation() {
            return this.channelUrl.trim() !== '' && 
                    (this.channelUrl.startsWith('https://t.me/') || 
                    this.channelUrl.startsWith('t.me/'));
        },
        isAuthenticated() {
            return !!localStorage.getItem('token');
        }
    },
    methods: {
        async addChannel() {
            if (!this.isAuthenticated) {
                this.showAlert = true;
                this.alertVariant = 'warning';
                this.alertMessage = 'Для добавления канала необходимо авторизоваться';
                return;
            }

            try {
                await axios.post('/api/protected/channels', {
                    channel_url: this.channelUrl
                });
                
                this.showAlert = true;
                this.alertVariant = 'success';
                this.alertMessage = 'Канал успешно добавлен';
                this.channelUrl = '';
                
                await this.fetchChannels();
                await this.fetchMessages();
            } catch (error) {
                this.showAlert = true;
                this.alertVariant = 'danger';
                this.alertMessage = 'Ошибка при добавлении канала: ' + 
                    (error.response?.data?.message || error.message);
            }
        },
        async removeChannel(channelUrl) {
            try {
                await axios.delete('/api/protected/channels', {
                    data: { channel_url: channelUrl }
                });
                
                this.showAlert = true;
                this.alertVariant = 'success';
                this.alertMessage = 'Канал успешно удален';
                
                await this.fetchChannels();
                await this.fetchMessages();
            } catch (error) {
                this.showAlert = true;
                this.alertVariant = 'danger';
                this.alertMessage = 'Ошибка при удалении канала: ' + 
                    (error.response?.data?.message || error.message);
            }
        },
        async fetchChannels() {
            if (!this.isAuthenticated) return;
            
            try {
                const response = await axios.get('/api/protected/channels');
                this.channels = response.data.channels || [];
            } catch (error) {
                console.error('Error fetching channels:', error);
            }
        },
        async fetchMessages() {
            if (!this.isAuthenticated) return;
            
            this.loading = true;
            try {
                const response = await axios.get('/api/protected/channels/messages');
                this.messages = response.data.messages || [];
            } catch (error) {
                console.error('Error fetching messages:', error);
                this.showAlert = true;
                this.alertVariant = 'danger';
                this.alertMessage = 'Ошибка при загрузке сообщений';
            } finally {
                this.loading = false;
            }
        },
        formatDate(dateStr) {
            if (!dateStr) return '';
            return new Date(dateStr).toLocaleString();
        },
        startAutoRefresh() {
            this.refreshInterval = setInterval(() => {
                this.fetchMessages();
            }, 60000); // Обновляем каждую минуту
        },
        stopAutoRefresh() {
            if (this.refreshInterval) {
                clearInterval(this.refreshInterval);
                this.refreshInterval = null;
            }
        }
    },
    async mounted() {
        if (this.isAuthenticated) {
            await this.fetchChannels();
            await this.fetchMessages();
            this.startAutoRefresh();
        }
    },
    beforeDestroy() {
        this.stopAutoRefresh();
    }
}
</script>

<style scoped>
.tg-news-container {
    width: 100%;
    max-width: 1500px;
    display: flex;
    flex-direction: column;
    align-items: start;
    margin: 0 auto;
    padding: 20px;
}

.tg-news-header {
    text-align: center;
    margin-bottom: 30px;
}

.channel-form {
    width: 100%;
}
.channel-form input {
    width: 90%;
}
.channel-form button {
    margin: 10px 20px;
    white-space: nowrap;
}
.messages-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
    margin-top: 20px;
}

.message-card {
    height: 100%;
    text-align: left;
    font-size: 15px;
}

.message-card .card {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.message-card .card-text {
    flex-grow: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 4;
    -webkit-box-orient: vertical;
}

.channels-list {
    width: 100%;
    margin: 0 auto;
}
</style>
