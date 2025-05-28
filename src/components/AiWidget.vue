<template>
    <div class="ai-widget" :class="{ expanded: isAIWidgetExpanded }">
        <div class="ai-widget-head" @click="toggleAIWidget">Задать вопрос ИИ</div>

        <div class="ai-widget-content">
            <p class="ai-widget-info">Задайте вопрос по непонятной вам теме, и искусственный интеллект объяснит вам её.</p>

            <!-- Conditional Loading and Response -->
            <p v-if="isAILoading">Загрузка...</p>
            <p v-else-if="aiResponse">{{ aiResponse }}</p>

            <label>
            <input type="text" placeholder="Введите ваш вопрос..." v-model="aiQuestion">
            <button @click="askAI" :disabled="isAILoading">></button>
            </label>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            isAIWidgetExpanded: false,
            aiQuestion: '',
            aiResponse: '',
            isAILoading: false // Add loading state
        }
    },
    methods: {
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
    }
}
</script>

<style scoped>
    .ai-widget {
        position: fixed;
        bottom: 20px;
        right: 20px;
        background-color: rgba(217, 217, 217, 1);
        color: white;
        border-radius: 10px;
        cursor: pointer;
        box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
        transition: background-color 0.3s ease;
        z-index: 1000;
        max-width: 20vw;
        min-width: 10vw;
    }
    .ai-widget-head{
        text-wrap: none;
        background-color: #7986cb;
        text-align: center;
        width: 100%;
        padding: 10px 5%;
        border-radius: 10px 10px 0 0;
    }
    .ai-widget-head:hover {
        background-color: #5e6fa3;
    }
    .ai-widget-expanded {
        background-color: rgba(217, 217, 217, 1);
        color: #333;
        padding: 20px;
        border-radius: 10px;
        min-width: 300px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    .ai-widget p {
        margin-bottom: 10px;
        color: #333;
    }
    .ai-widget input[type="text"] {
        width: 70%;
        padding: 5px;
        border: 1px solid #aaa;
        border-radius: 5px;
        margin-bottom: 10px;
    }
    .ai-widget-info{
        font-style: italic;
        font-size: 10pt;
    }
    .ai-widget button {
        background-color: #7986cb;
        color: white;
        border: none;
        padding: 5px 9px;
        border-radius: 5px;
        cursor: pointer;
        transition: background-color 0.3s ease;
        margin-left: 10px;
    }
    .ai-widget button:hover {
        background-color: #5e6fa3;
    }
    .ai-widget-content {
        display: none;
        background-color: rgba(217, 217, 217, 1);
        margin: 10px;
    }
    .ai-widget.expanded .ai-widget-content {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        min-height: 24vh;
    }
    @media (max-width: 400px) {
        .ai-widget-expanded {
        width: 90%;
        }
    }
</style>