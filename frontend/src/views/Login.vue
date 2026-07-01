<template>
  <div class="loginpage">
    <!-- Animated Background Particles -->
    <div class="bg-particles">
      <div class="particle p1"></div>
      <div class="particle p2"></div>
      <div class="particle p3"></div>
      <div class="particle p4"></div>
      <div class="particle p5"></div>
      <div class="particle p6"></div>
      <div class="particle p7"></div>
      <div class="particle p8"></div>
    </div>

    <div class="auth-section">
      <LoginForm />
    </div>
  </div>
</template>

<script setup>
import '../assets/css/login.css'
import LoginForm from '../components/auth/LoginForm.vue'
</script>
