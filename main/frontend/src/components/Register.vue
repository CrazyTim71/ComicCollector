<template>
    <div class="vh-100 d-flex flex-column justify-content-center align-items-center gradient">
        <div class="text-center mb-5">
            <img src="/favicon.ico" alt="Logo" onclick="window.location.href='/'" class="mb-2" style="width: 100px; height: auto; ">
            <h2 class="text-white">Comic Collector</h2>
        </div>

        <div id="app" class="p-5 rounded shadow bg-white" style="width: 450px;">
            <div class="text-center mb-4">
                <h1 class="mb-3">Create Account</h1>
            </div>
            <div class="form-floating mb-3">
                <input v-model="username" type="text" class="form-control" id="floatingInput" placeholder="Username">
                <label for="floatingInput">Username</label>
            </div>
            <div class="form-floating mb-3">
                <input v-model="password" type="password" class="form-control" id="floatingPassword" placeholder="Password" @keyup.enter="register()">
                <label for="floatingPassword">Password</label>
            </div>
            <div class="form-floating mb-3">
                <input v-model="passwordRepeated" type="password" class="form-control" id="floatingPassword" placeholder="Repeat Password" @keyup.enter="register()">
                <label for="floatingPassword">Repeat Password</label>
            </div>
            <div class="d-flex justify-content-center">
                <button @click="register()" class="btn btn-primary w-100 btn-custom">Submit</button>
            </div>
            <div v-if="errorMsg" class="alert alert-danger text-center mt-3 mb-0" role="alert">{{ errorMsg }}</div>
        </div>
    </div>
</template>
    
<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useHead } from '@vueuse/head';
import axios from 'axios';

export default defineComponent({
  name: 'Login',
  setup() {
    const username = ref('')
    const password = ref('')
    const passwordRepeated = ref('')
    const errorMsg = ref('')

    useHead({
        title: 'Register | Comic Collector',
    });

    const register = async () => {
        errorMsg.value = ''

        if (username.value === '' || password.value === '' || passwordRepeated.value === '') {
            errorMsg.value = 'Please fill in all fields'
            return
        }

        if (password.value !== passwordRepeated.value) {
            errorMsg.value = 'Passwords do not match'
            return
        }

        try {
            let response = await axios.post('/api/v1/register', {
                username: username.value,
                password: password.value,
                passwordRepeated: passwordRepeated.value
            })
            
            if (response.status === 302 || response.status === 303) {
                window.location.href = response.headers['location'] || '/';
            } else if (response.status >= 200 && response.status < 300) {
                window.location.href = response.data.url || '/';
            } else {
                errorMsg.value = response.data.msg;
            }
        } catch (error) {
            errorMsg.value = 'Registration failed'
        }
    }

    return {
        username,
        password,
        passwordRepeated,
        errorMsg,
        register
    }
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