const {
  createApp,
  ref,
  onMounted,
  onBeforeUnmount
} = Vue

createApp({
  data() {
    return {
      newsData: [],
      currentCategory: 'all',
      currentDate: new Date(),
      offset: 0,
      limit: 15,
      hasMore: true,
      isMobile: false,
      menuVisible: false,
      isAIWidgetExpanded: false,
      aiQuestion: '',
      aiResponse: '',
      isAILoading: false // Add loading state
    };
  },
  computed: {
    filteredNews() {
      return this.currentCategory === 'all' ?
        this.newsData :
        this.newsData.filter(newsItem => newsItem.tags === this.currentCategory);
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
    fetchNewsAll() {
      axios.get(`/news?offset=15`)
        .then(response => {
          if (response.data.length < this.limit) {
            this.hasMore = false;
          }
          this.newsData = [...this.newsData, ...response.data];
        })
        .catch(err => console.error(err));
    },
    fetchNews() {
      let url = `/news?offset=${this.offset}`;
      url += `&category=${this.currentCategory}`;

      axios.get(url)
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
    toggleAIWidget() {
      this.isAIWidgetExpanded = !this.isAIWidgetExpanded;
    },
    async askAI() {
      this.isAILoading = true;
      this.aiResponse = ''; // Clear previous response
      try {
        const response = await axios.post('/ask', {
          prompt: this.aiQuestion
        });
        this.aiResponse = response.data.content; // Get the "answer" property
      } catch (error) {
        console.error('Error asking AI:', error);
        this.aiResponse = 'Произошла ошибка при отправке запроса.';
      } finally {
        this.isAILoading = false;
      }
    }
  },
  mounted() {
    this.fetchNewsAll();
    this.intervalId = setInterval(() => {
      this.currentDate = new Date();
    }, 60000);
    this.checkScreenSize();
    window.addEventListener('resize', this.checkScreenSize);
  },
  beforeUnmount() {
    clearInterval(this.intervalId);
    window.removeEventListener('resize', this.checkScreenSize);
  },
}).mount('#app');