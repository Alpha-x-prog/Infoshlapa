<template>
    <div class="date">
        <span id="current-date">{{ formattedDateTime }}</span>
        <hr>
    </div>
</template>

<script>
export default {
    data() {
        return {
            currentDate: new Date(),
        }
    },
    computed: {
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
    mounted() {
        this.intervalId = setInterval(() => {
            this.currentDate = new Date();
        }, 60000)
    },
    beforeDestroy() {
        clearInterval(this.intervalId);
    }
}
</script>

<style scoped>
    .date{
        width: 100%;
        margin: 3% 0 ;
        color: rgba(69, 78, 82, 1);
        font-style: italic;
    }
    .date hr{
        color: rgba(69, 78, 82, 1);
        border: none;
        border-top: 1px solid rgba(69, 78, 82, 1);
    }
    @media (max-width: 700px) {
        .date{
            font-size: 0.9rem;
        }
    }
</style>