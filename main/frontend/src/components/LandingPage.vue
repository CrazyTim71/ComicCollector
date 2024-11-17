<template>
    <div class="vh-100 d-flex flex-column justify-content-center align-items-center gradient">
        <div class="text-center mb-4">
            <img src="/favicon.ico" alt="Logo" class="mb-2" style="width: 100px; height: auto;">
        </div>

        <!-- Hero Section -->
        <div class="hero-section text-center text-white">
                <h1 class="fs-1">Welcome to Comic Collector</h1>
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

<script>
import axios from 'axios';
import.meta.env.VUE_APP_API_BASE_URL
import '../css/LandingPage.css'

export default {
    name: 'LandingPage',
    data() {
        return {
            signupEnabled: false
        };
    },
    mounted() {
        axios.get('/api/v1/register/check')
            .then(response => {
                this.signupEnabled = response.data.signupEnabled;
            })
            .catch(error => {
                console.error('Error checking signup status:', error);
            });
    }
};
</script>
