<template>   
    <div class="vh-100 d-flex flex-column justify-content-center align-items-center gradient">
        <div class="text-center mb-4">
            <img src="/favicon.ico" alt="Logo" class="mb-2" style="width: 100px; height: auto;">
        </div>

        <!-- Hero Section -->
        <div class="hero-section text-center text-white">
                <h1 class="fs-1">Welcome to {{ title }}</h1>
                <p class="hero-description fs-5 mb-4 fw-light"><em>Your ultimate app for managing comic book collections</em></p>
                <a href="/login" class="btn btn-primary btn-lg mt-4">Login</a>
            </div>

        <!-- Footer Links -->
        <div class="footer-links text-center mt-5 text-white">

            <a v-if="signupEnabled" href="/register" class="text-white me-3 ms-4">Create an Account</a>
            <a href="/privacy" class="text-white me-3 ms-4">Privacy Policy</a>
            <a href="/terms" class="text-white ms-4">Terms of Service</a>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { useHead } from '@vueuse/head';
import axios from 'axios';

export default defineComponent({
    name: 'LandingPage',
    setup() {
        useHead({
            title: 'Comic Collector',
            meta: [
                {
                    name: 'description',
                    content: 'Your ultimate app for managing comic book collections'
                }
            ]
        });
    },
    data() {
        return {
            signupEnabled: false,
            title: 'Comic Collector'
        };
    },
    methods: {
        checkSignup() {
            axios.get('/api/v1/register/check')
                .then(response => {
                    this.signupEnabled = response.data.signupEnabled;
                })
                .catch(error => {
                    console.error('Error checking signup status:', error);
                });
        }
    },
    mounted() {
        this.checkSignup();
    }
});
</script>

<style scoped>
/* Background Gradient */
.gradient {
    background: linear-gradient(90deg, rgba(2,0,36,1) 0%, rgba(233,216,197,1) 0%, rgba(139,201,204,1) 46%, rgba(10,180,214,1) 100%);
}

/* Custom Button Styles */
.btn-primary {
    background-color: #007bff;
    border-color: #007bff;
    transition: background-color 0.3s, box-shadow 0.3s;
}

.btn-primary:hover {
    background-color: #0056b3;
    box-shadow: 0 0 10px rgba(0, 91, 187, 0.5);
}

/* Hero Description */
.hero-description {
    color: #f5f5f5;
}

/* Footer Links */
.footer-links a {
    font-size: 0.9rem;
    text-decoration: underline;
    transition: color 0.2s ease;
}
</style>
