const { createApp } = Vue;

createApp({
    data() {
        return {
            newsList: [],
            offset: 0,
            limit: 15,
            hasMore: true
        };
    },
    methods: {
        fetchNews() {
            axios.get(`/news?offset=${this.offset}`)
                .then(response => {
                    if (response.data.length < this.limit) {
                        this.hasMore = false;
                    }
                    this.newsList = [...this.newsList, ...response.data];
                    this.offset += this.limit;
                })
                .catch(err => console.error(err));
        },
        loadMore() {
            this.fetchNews();
        },
        formatDate(date) {
            return new Date(date).toLocaleString();
        }
    },
    mounted() {
        this.fetchNews();
    }
}).mount("#app");
