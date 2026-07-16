<template>
  <div class="loginpage" style="display: flex; min-height: 100vh; width: 100%; background-color: #000000; color: #E2E8F0; font-family: 'Outfit', sans-serif; position: relative; overflow: hidden; align-items: stretch; justify-content: center; padding: 0;">
    <!-- Animated Background Particles -->
    <div class="bg-particles" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; overflow: hidden; z-index: 0; pointer-events: none; display: block;">
      <div class="particle p1"></div><div class="particle p2"></div><div class="particle p3"></div><div class="particle p4"></div>
      <div class="particle p5"></div><div class="particle p6"></div><div class="particle p7"></div><div class="particle p8"></div>
    </div>

    <!-- LEFT PANEL -->
    <div class="left-panel" style="flex: 0 0 470px; background: transparent; padding: 60px 80px 60px 80px; display: flex; flex-direction: column; position: relative; z-index: 2; border-right: 1px solid rgba(255, 255, 255, 0.05);">
      <div class="brand-logo" style="font-size: 28px; font-weight: 800; color: #FFFFFF; display: flex; align-items: center; gap: 12px; margin-bottom: 2rem; justify-content: flex-start; width: 100%;">
        <img src="../assets/icons/favicon.svg" alt="MoneymakePro" style="width: 28px; height: 28px;" />
        MoneymakePro
      </div>
      
      <div class="marketing-content" style="flex: 1; display: flex; flex-direction: column; justify-content: center;">
        <div style="text-align: left; margin-bottom: 2rem;">
          <img src="../assets/images/wallet-illustration.png" alt="Wallet Illustration" style="max-width: 100%; mix-blend-mode: screen;" />
        </div>
        
        <h2 class="marketing-headline" style="font-size: 32px; font-weight: 700; line-height: 1.2; margin-bottom: 1rem; color: #FFFFFF;">
          Welcome back to<br/>
          <span style="color: #A855F7;">MoneymakePro</span>
        </h2>
        <p class="marketing-subheadline" style="font-size: 15px; color: #94a3b8; margin-bottom: 2.5rem; line-height: 1.5;">
          Sign in to access your dashboard, monitor your investments, and grow your wealth.
        </p>
      </div>
    </div>

    <!-- RIGHT PANEL -->
    <div class="right-panel" style="flex: 0 0 560px; position: relative; display: flex; align-items: center; justify-content: center; padding: 40px; z-index: 2;">
      <div class="auth-section" style="width: 100%; max-width: 440px;">
        <LoginForm />
      </div>
    </div>
  </div>
</template>

<script setup>
import '../assets/css/login.css'
import LoginForm from '../components/auth/LoginForm.vue'
</script>

<style scoped>
@media (max-width: 1024px) {
  .left-panel {
    display: none !important;
  }
  .right-panel {
    padding: 24px !important;
  }
}
</style>

