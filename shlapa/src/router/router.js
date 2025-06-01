import { createRouter, createWebHistory } from 'vue-router'
import Main from '@/pages/main'
import TgNews from '@/pages/tgnews'
import Profile from '@/pages/profile'
import Bookmarks from '@/pages/bookmarks'

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
    },
    {
        path: '/bookmarks',
        component: Bookmarks,
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    routes,
    history: createWebHistory(process.env.BASE_URL)
})

// Navigation guard to check authentication
router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (!localStorage.getItem('token')) {
            next('/profile')
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router;
