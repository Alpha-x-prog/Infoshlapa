<template>
    <div class="news-card" @click="toggleContent">
        <div class="news-image-container" :class="{ hidden: showText }">
            <img :src="imageUrl" alt="News Image" class="news-image">
            <span class="news-category">{{ news.tags || 'Без категории' }}</span>
        </div>
        <div class="news-content" @click="toggleContent">
            <h3 class="news-title">{{ news.title || 'Без заголовка' }}</h3>
            <p class="news-text" v-if="showText">{{ news.description || 'Нет описания' }}</p>
        </div>
        <div class="news-footer">
            <span class="news-date">{{ formatDate(news.publishedAt) }}</span>
            <div class="news-actions">
                <button 
                    v-if="isAuthenticated"
                    class="bookmark-btn"
                    @click.stop="toggleBookmark"
                    :title="isBookmarked ? 'Удалить из закладок' : 'Добавить в закладки'"
                >
                    <i :class="['fas', isBookmarked ? 'fa-bookmark' : 'fa-bookmark-o']"></i>
                    закладка
                </button>
                <a :href="news.url || '#'" class="news-source" target="_blank">Источник...</a>
            </div>
        </div>
    </div>
</template>

<script>
import axios from '@/utils/axios'

export default {
    props: {
        news: {
            type: Object,
            required: true,
            default: () => ({
                title: '',
                description: '',
                url: '#',
                urlToImage: '',
                publishedAt: new Date(),
                tags: 'Без категории'
            })
        }
    },
    data() {
        return {
            showText: false,
            isBookmarked: false
        };
    },
    computed: {
        imageUrl() {
            return this.news.urlToImage || 'https://via.placeholder.com/350x250?text=No+Image';
        },
        isAuthenticated() {
            return !!localStorage.getItem('token');
        }
    },
    methods: {
        toggleContent() {
            this.showText = !this.showText;
        },
        formatDate(date) {
            if (!date) return 'Дата неизвестна';
            return new Date(date).toLocaleString();
        },
        async toggleBookmark() {
            if (!this.isAuthenticated) {
                this.$router.push('/profile');
                return;
            }

            try {
                if (this.isBookmarked) {
                    await axios.delete('/api/protected/bookmarks', {
                        data: { 
                            newsId: this.news.article_id
                        }
                    });
                    this.isBookmarked = false;
                } else {
                    await axios.post('/api/protected/bookmarks', {
                        newsId: this.news.article_id
                    });
                    this.isBookmarked = true;
                }
            } catch (error) {
                console.error('Error toggling bookmark:', error);
            }
        }
    },
    async mounted() {
        if (this.isAuthenticated) {
            try {
                const response = await axios.get('/api/protected/bookmarks');
                if (response.data && Array.isArray(response.data)) {
                    this.isBookmarked = response.data.some(bookmark => bookmark.article_id === this.news.article_id);
                } else {
                    this.isBookmarked = false;
                }
            } catch (error) {
                console.error('Error checking bookmark status:', error);
                this.isBookmarked = false;
            }
        }
    }
}
</script>

<style scoped>
    .news-card {
        margin: 20px;
        width: 350px;
        min-height: 400px;
        display: flex;
        flex-direction: column;
        align-items: start;
        justify-content: space-between;
        background-color: #fff;
        border-radius: 8px;
        box-shadow: 0 3px 5px rgba(0, 0, 0, 0.1);
        overflow: hidden; 
        cursor: pointer; 
        transition: background-color 0.3s ease;
    }
    .news-card:hover {
        background-color: rgb(240, 240, 240); /* Lighter background on hover */
    }
    .news-image-container {
        position: relative;
        height: 250px;
        overflow: hidden;
        transition: opacity 0.3s ease; /* Fade effect on image */
    }
    .news-image-container.hidden {
        opacity: 0;
        height: 0; /* Collapse the image container */
    }
    .news-image {
        width: 100%;
        height: 100%;
        object-fit: cover;
        display: block; /* Corrects possible image rendering issues */
    }
    .news-category {
        position: absolute;
        top: 10px;
        left: 10px;
        background-color: #2e2e2ea4;
        color: aliceblue;
        padding: 5px 10px;
        border-radius: 5px;
        font-size: 12pt;
    }
    .news-content {
        padding: 15px;
        text-align: left;
        transition: padding 0.3s ease;
        cursor: pointer;
    }
    .news-title {
        font-size: 18pt;
        margin-bottom: 10px;
        line-height: 1.0;
        color: #333;
    }
    .news-text {
        font-size: 15pt;
        color: #555;
        line-height: 1.2;
        margin-bottom: 10px;
        overflow: hidden; /* Prevents content from overflowing */
        transition: height 0.9s ease, opacity 0.9s ease;
        height: auto;
    }
    .news-footer {
        width: 90%;
        display: flex;
        justify-content: space-between; 
        align-items: center;
        font-size: 12pt;
        color: #777;
        margin: 0 5% 15px 5%; 
    }
    .news-date {
        font-style: italic; 
    }
    .news-actions {
        display: flex;
        align-items: center;
        gap: 15px;
    }
    .news-source {
        text-decoration: none;
        color: #549cd6; 
        font-weight: bold; 
    }
    .news-source:hover{
        color: #91b5d3;
    }
    .bookmark-btn {
        background: none;
        border: none;
        color: #549cd6;
        cursor: pointer;
        padding: 5px;
        font-size: 1.2em;
        transition: color 0.3s ease;
    }
    .bookmark-btn:hover {
        color: #91b5d3;
    }
    .bookmark-btn i {
        pointer-events: none;
    }
</style>
