import { createRouter, createWebHistory } from 'vue-router'
import Main from '@/pages/main'
import TgNews from '@/pages/tgnews'
import Profile from '@/pages/profile'

const routes = [
    {
        path: '/',
        component: Main
    },
    {
        path: '/tgnews',
        component: TgNews
    },
    {
        path: '/profile',
        component: Profile
    }
]

const router = createRouter({
    routes,
    history: createWebHistory(process.env.BASE_URL)
})

export default router;
