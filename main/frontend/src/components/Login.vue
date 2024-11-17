<template>
<div class="container vw-100 vh-100">
    <div class="row justify-content-center">
    <div class="col-md-6">
        <div class="card">
        <div class="card-header">
            <h3 class="text-center">Login</h3>
        </div>
        <div class="card-body">
            <form @submit.prevent="login">
            <div class="form-group mb-3">
                <label for="email">Email</label>
                <input type="email" v-model="email" class="form-control" id="email" required />
            </div>
            <div class="form-group mb-3">
                <label for="password">Password</label>
                <input type="password" v-model="password" class="form-control" id="password" required />
            </div>
            <button type="submit" class="btn btn-primary w-100">Login</button>
            </form>
        </div>
        </div>
    </div>
    </div>
</div>
</template>

<script setup lang="ts">
    import { ref } from 'vue';

    const email = ref('');
    const password = ref('');

    const login = () => {
        console.log('Email:', email.value);
        console.log('Password:', password.value);
        // Add your login logic here
    };
</script>