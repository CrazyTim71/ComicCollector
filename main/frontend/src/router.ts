import { createRouter, createWebHistory } from 'vue-router';
import Login from './components/Login.vue';
import LandingPage from './components/LandingPage.vue';
import Register from './components/Register.vue';

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
    {
        path: '/register',
        name: 'Register',
        component: Register,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;