<template>
    <div class="bookmarks-container">
        <div class="bookmarks-header">
            <p v-if="!isAuthenticated" class="unauthorized-message">
                Для использования этой функции необходимо авторизоваться
            </p>
            <p v-else-if="bookmarks.length === 0" class="no-bookmarks">
                У вас пока нет сохраненных новостей. 
                Нажмите на иконку закладки <i class="bi bi-bookmark"></i> рядом с новостью, чтобы сохранить ее.
            </p>
        </div>
        
        <div class="bookmarks-list" v-if="isAuthenticated && bookmarks.length > 0">
            <NewsCard
                v-for="news in bookmarks"
                :key="news.id"
                :news="news"
                @bookmark-updated="handleBookmarkUpdate"
            />
        </div>

        <div v-if="loading" class="loading-spinner">
            <b-spinner label="Loading..."></b-spinner>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import NewsCard from './NewsCard.vue';

export default {
    name: 'BookmarksContent',
    components: {
        NewsCard
    },
    data() {
        return {
            bookmarks: [],
            loading: false
        }
    },
    computed: {
        isAuthenticated() {
            return localStorage.getItem('token') !== null;
        }
    },
    methods: {
        async fetchBookmarks() {
            if (!this.isAuthenticated) {
                return;
            }

            this.loading = true;
            try {
                const response = await axios.get('/api/bookmarks', {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });

                if (response.data.status === 'success') {
                    this.bookmarks = response.data.bookmarks;
                }
            } catch (error) {
                console.error('Error fetching bookmarks:', error);
            } finally {
                this.loading = false;
            }
        },
        handleBookmarkUpdate({ id, isBookmarked }) {
            if (!isBookmarked) {
                this.bookmarks = this.bookmarks.filter(news => news.id !== id);
            }
        }
    },
    mounted() {
        this.fetchBookmarks();
    }
}
</script>

<style scoped>
.bookmarks-container {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
}

.bookmarks-header {
    margin-bottom: 30px;
    text-align: center;
}

.bookmarks-header h1 {
    color: #2c3e50;
    margin-bottom: 15px;
}

.no-bookmarks {
    color: #666;
    font-size: 1.1rem;
}

.loading-spinner {
    display: flex;
    justify-content: center;
    padding: 40px;
}

.bookmarks-list {
    display: grid;
    gap: 20px;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    justify-items: center;
}

.unauthorized-message {
    padding: 20px;
    background-color: #eff7c4;
    border-radius: 10px;
    border: 1px solid #636224;
    color: #636224;
    font-size: 1.1rem;
    text-align: center;
    margin: 2rem 0;
}
</style>
        