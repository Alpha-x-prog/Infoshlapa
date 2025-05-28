<template>
    <div class="news-section">
        <navbar @category-changed="setCategory" />
        <div class="news-container">
            <div v-if="loading" class="loading">Loading...</div>
            <div v-if="error" class="error">&#128711; {{ error }}</div>
            <template v-if="!loading && !error">
                <news-card v-for="newsItem in filteredNews" 
                          :key="newsItem.id" 
                          :news="newsItem">
                </news-card>
                <div v-if="hasMore" class="load-more">
                    <button @click="loadMore" :disabled="loading">
                        Загрузить еще
                    </button>
                </div>
            </template>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import NewsCard from './NewsCard.vue';
import Navbar from './Navbar.vue';

export default {
    name: 'News',
    components: {
        NewsCard,
        Navbar
    },
    data() {
        return {
            newsData: [],
            currentCategory: 'all',
            offset: 0,
            limit: 15,
            hasMore: true,
            loading: false,
            error: null
        };
    },
    computed: {
        filteredNews() {
            return this.currentCategory === 'all'
                ? this.newsData
                : this.newsData.filter(newsItem => newsItem.tags === this.currentCategory);
        }
    },
    methods: {
        setCategory(category) {
            this.currentCategory = category;
            this.newsData = [];
            this.offset = 0;
            this.hasMore = true;
            this.fetchNews();
        },
        async fetchNews() {
            try {
                this.loading = true;
                this.error = null;
                
                const url = `/news?offset=${this.offset}&limit=${this.limit}${
                    this.currentCategory !== 'all' ? `&category=${this.currentCategory}` : ''
                }`;
                
                const response = await axios.get(url);
                const newNews = response.data;
                
                if (newNews.length < this.limit) {
                    this.hasMore = false;
                }
                
                this.newsData = [...this.newsData, ...newNews];
                this.offset += this.limit;
            } catch (err) {
                this.error = 'Failed to load news. Please try again later.';
                console.error('Error fetching news:', err);
            } finally {
                this.loading = false;
            }
        },
        loadMore() {
            if (!this.loading && this.hasMore) {
                this.fetchNews();
            }
        }
    },
    mounted() {
        this.fetchNews();
    }
};
</script>

<style scoped>
.news-section {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
}

.news-container {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    user-select: none;
}

.loading, .error {
    width: 100%;
    text-align: center;
    padding: 20px;
    font-size: 18px;
}

.error { 
    color: #ac3a3a;
}

.load-more {
    width: 100%;
    text-align: center;
    margin: 20px 0;
}

.load-more button {
    padding: 10px 20px;
    background-color: #549cd6;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.3s;
}

.load-more button:hover {
    background-color: #91b5d3;
}

.load-more button:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
}
</style>