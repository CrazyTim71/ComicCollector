<template>
  <div class="vh-100 d-flex flex-column justify-content-center align-items-center gradient">
    <div class="text-center mb-5">
        <img src="/favicon.ico" alt="Logo" onclick="window.location.href='/'" class="mb-2" style="width: 100px; height: auto; ">
        <h2 class="text-white">Comic Collector</h2>
    </div>

    <div id="app" class="p-5 rounded shadow bg-white" style="width: 450px;">
      <div class="text-center mb-4">
          <h1 class="mb-3">Login</h1>
      </div>
      <div class="form-floating mb-3">
          <input v-model="username" type="text" class="form-control" id="floatingInput" placeholder="Username">
          <label for="floatingInput">Username</label>
      </div>
      <div class="form-floating mb-3">
          <input v-model="password" type="password" class="form-control" id="floatingPassword" placeholder="Password" @keyup.enter="login()">
          <label for="floatingPassword">Password</label>
      </div>
      <div class="d-flex justify-content-center">
          <button @click="login()" class="btn btn-primary w-100 btn-custom">Submit</button>
      </div>

      <div v-if="errorMsg" class="alert alert-danger text-center mt-3 mb-0" role="alert">{{ errorMsg }}</div>
    </div>

    <div v-if="signupEnabled" class="text-center mt-3">
      <a href="/register" class="text-white">Create an Account</a>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue';
import { useHead } from '@vueuse/head';
import axios from 'axios';

export default defineComponent({
  name: 'Login',
  setup() {
    const username = ref('');
    const password = ref('');
    const errorMsg = ref('');
    const signupEnabled = ref(false);

    useHead({
        title: 'Login | Comic Collector',
    });

    const checkSignup = () => {
      axios.get('/api/v1/register/check')
        .then(response => {
          signupEnabled.value = response.data.signupEnabled;
        })
        .catch(error => {
          console.error('Error checking signup status:', error);
        });
    };

    onMounted(() => {
      checkSignup();
    });

    const login = async () => {
      errorMsg.value = '';

    if (username.value === "" || password.value === "") {
      errorMsg.value = 'Please fill in all fields';
      return;
    }

    try {
        let response = await fetch('/api/v1/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: username.value,
                password: password.value
            })
        });

        if (response.status === 302 || response.status === 303) {
            window.location.href = response.headers.get('Location') || '/';
        } else if (response.ok) {
            window.location.href = response.url;
        } else {
            const result = await response.json();
            errorMsg.value = result.msg;

            password.value = "";
            return;
        }
    } catch (exception) {
      console.error('An unexpected error occurred:', exception);
      errorMsg.value = 'An unexpected error occurred.';
    }
    };

    return {
      username,
      password,
      errorMsg,
      signupEnabled,
      checkSignup,
      login
    };
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