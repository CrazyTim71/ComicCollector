import { createRouter, createWebHistory } from 'vue-router';
import Login from './components/Login.vue';
import LandingPage from './components/LandingPage.vue';

const routes = [
    {
        path: '/',
        name: 'Comic Collector',
        component: LandingPage,
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;