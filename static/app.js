const { createApp, ref, onMounted, onBeforeUnmount } = Vue

createApp({
    data() {
        return {
            newsData: [],
            currentCategory: 'all',
            currentDate: new Date(),
            offset: 0,
            limit: 15,
            hasMore: true,
            isMobile: false, // Add isMobile to data
            menuVisible: false,
        };
    },
    computed: {
        filteredNews() {
            return this.currentCategory === 'all'
                ? this.newsData
                : this.newsData.filter(newsItem => newsItem.tags === this.currentCategory);
        },
        formattedDateTime() {
            const months = [
                "января", "февраля", "марта", "апреля", "мая", "июня",
                "июля", "августа", "сентября", "октября", "ноября", "декабря"
            ];
            const day = this.currentDate.getDate();
            const monthIndex = this.currentDate.getMonth();
            const year = this.currentDate.getFullYear();
            const hours = String(this.currentDate.getHours()).padStart(2, '0');
            const minutes = String(this.currentDate.getMinutes()).padStart(2, '0');
            return `${day} ${months[monthIndex]} ${year}, ${hours}:${minutes}`;
        }
    },
    components: {
        'news-card': {
            props: ['news'],
            template: `
                <div class="news-card" @click="toggleContent">
                    <div class="news-image-container" :class="{ hidden: showText }">
                        <img :src="news.urlToImage" alt="News Image" class="news-image">
                        <span class="news-category">{{ news.tags }}</span>
                    </div>
                    <div class="news-content">
                        <h3 class="news-title">{{ news.title }}</h3>
                        <p class="news-text" v-if="showText">{{ news.description }}</p>
                    </div>
                    <div class="news-footer">
                        <span class="news-date">{{ formatDate(news.publishedAt) }}</span>
                        <a :href="news.url" class="news-source">Источник...</a>
                    </div>
                </div>
            `,
            data() {
                return {
                    showText: false
                };
            },
            methods: {
                toggleContent() {
                    this.showText = !this.showText;
                },
                formatDate(date) {
                    return new Date(date).toLocaleString();
                }
            }
        }
    },
    methods: {
        setCategory(tags) {
            this.currentCategory = tags;
        },
        fetchNews() {
            axios.get(`/news?offset=${this.offset}`)
                .then(response => {
                    if (response.data.length < this.limit) {
                        this.hasMore = false;
                    }
                    this.newsData = [...this.newsData, ...response.data];
                    this.offset += this.limit;
                })
                .catch(err => console.error(err));
        },
        loadMore() {
            this.fetchNews();
        },
        checkScreenSize() {
          this.isMobile = window.innerWidth <= 700;
        },
        toggleMenu() {
          this.menuVisible = !this.menuVisible;
        },
    },
    mounted() {
      this.fetchNews();
      this.intervalId = setInterval(() => {
        this.currentDate = new Date();
      }, 60000);
      this.checkScreenSize(); // Call on mount
      window.addEventListener('resize', this.checkScreenSize); // And on resize
    },
    beforeUnmount() {
        clearInterval(this.intervalId);
        window.removeEventListener('resize', this.checkScreenSize);
    },
}).mount('#app');