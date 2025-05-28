import { createApp } from 'vue'
import App from './App'
import router from '@/router/router'
import axios from 'axios'

import BootstrapVue3 from 'bootstrap-vue-3'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue-3/dist/bootstrap-vue-3.css'

// Настройка axios
axios.defaults.baseURL = '/api';

const app = createApp(App)

app
    .use(BootstrapVue3)
    .use(router)
    .mount('#app');
