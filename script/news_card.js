const { createApp } = Vue

createApp({
data() {
    return {
      newsData: [] // Initialize as an empty array
    };
},
components: {
    'news-card': {
      props: ['news'], // Receive the news item as a prop
        template: `
        <div class="news-card" @click="toggleContent">
            <div class="news-image-container" :class="{ hidden: showText }">
                <img :src="news.imageUrl" alt="News Image" class="news-image">
                <span class="news-category">{{ news.category }}</span>
            </div>
            <div class="news-content">
                <h3 class="news-title">{{ news.title }}</h3>
                <p class="news-text" v-if="showText">{{ news.text }}</p>
            </div>
            <div class="news-footer">
                <span class="news-date">{{ news.date }}</span>
                <a :href="news.sourceUrl" class="news-source">Источник...</a>
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
        }
        }
    }
},
mounted() {
    fetch('../json/news.json')
        .then(response => response.json())
        .then(data => {
        this.newsData = data; // Assign the entire array to newsData
        })
        .catch(error => console.error('Error fetching data:', error));
    }
}).mount('#app')