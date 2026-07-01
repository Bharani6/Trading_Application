<template>
  <div class="error-card">
    <svg class="error-alert-icon" viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
      <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
      <line x1="12" y1="9" x2="12" y2="13"></line>
      <line x1="12" y1="17" x2="12.01" y2="17"></line>
    </svg>
    <div class="error-content">
      <h3 class="error-title">{{ title }}</h3>
      <p class="error-message">{{ message }}</p>
    </div>
    <button v-if="retry" class="retry-btn" @click="$emit('retry')">
      <svg class="retry-icon-svg" viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="23 4 23 10 17 10"></polyline>
        <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
      </svg> Retry
    </button>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: 'An error occurred'
  },
  message: {
    type: String,
    required: true
  },
  retry: {
    type: Boolean,
    default: false
  }
})
defineEmits(['retry'])
</script>

<style scoped>
.error-card {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  background: rgba(244, 63, 94, 0.1);
  border: 1px solid rgba(244, 63, 94, 0.2);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
}
.error-alert-icon {
  color: #f43f5e;
  margin-top: 2px;
}
.error-content {
  flex: 1;
}
.error-title {
  color: #f43f5e;
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 4px 0;
}
.error-message {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin: 0;
}
.retry-btn {
  background: transparent;
  border: 1px solid rgba(244, 63, 94, 0.3);
  color: #f43f5e;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}
.retry-btn:hover {
  background: rgba(244, 63, 94, 0.1);
}
.retry-icon-svg {
  margin-right: 4px;
  vertical-align: middle;
}
</style>
