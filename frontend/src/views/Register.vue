<template>
  <div class="register-split-page">
    <!-- Animated Background Particles -->
    <div class="bg-particles">
      <div class="particle p1"></div><div class="particle p2"></div><div class="particle p3"></div><div class="particle p4"></div>
      <div class="particle p5"></div><div class="particle p6"></div><div class="particle p7"></div><div class="particle p8"></div>
    </div>

    <div class="right-panel">
      <RegisterForm />
    </div>
  </div>
</template>

<script setup>
import '../assets/css/register.css'
import RegisterForm from '../components/auth/RegisterForm.vue'
</script>
