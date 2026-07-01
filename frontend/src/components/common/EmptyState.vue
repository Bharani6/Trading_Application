<template>
  <div class="empty-state">
    <i :class="icon" class="empty-icon"></i>
    <p class="empty-message">{{ message }}</p>
    <slot></slot>
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: String,
    default: 'No data available'
  },
  icon: {
    type: String,
    default: 'fas fa-inbox'
  }
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 1rem;
  text-align: center;
  color: var(--text-muted);
}
.empty-icon {
  font-size: 32px;
  margin-bottom: 1rem;
  color: var(--border-color);
}
.empty-message {
  font-size: 15px;
  font-weight: 500;
}
</style>
