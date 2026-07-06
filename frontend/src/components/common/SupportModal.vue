<template>
  <div class="modal-overlay" v-if="isOpen" @click.self="close">
    <div class="modal-content" v-bind:class="{ 'is-closing': isClosing }">
      <div class="modal-header">
        <h3>Contact Support</h3>
        <button class="close-btn" @click="close">&times;</button>
      </div>

      <div class="modal-body" v-if="!submitted">
      <p class="subtitle">Need Assistance? We're Here to Help</p>
        
        <form @submit.prevent="submitForm" class="support-form">
          <div class="form-group">
            <label for="name">Name</label>
            <input 
              type="text" 
              id="name" 
              v-model="form.name" 
              required
              placeholder="Your full name"
            />
          </div>

          <div class="form-group">
            <label for="email">Email</label>
            <input 
              type="email" 
              id="email" 
              v-model="form.email" 
              required
              placeholder="your@email.com"
            />
          </div>

          <div class="form-group">
            <label for="message">Message</label>
            <textarea 
              id="message" 
              v-model="form.message" 
              required
              placeholder="Describe your issue or question..."
              rows="4"
            ></textarea>
          </div>

          <div class="error-message" v-if="error">{{ error }}</div>

          <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
            {{ loading ? 'Sending...' : 'Send Message' }}
          </button>
        </form>
      </div>
      
      <div class="modal-body success-state" v-else>
        <div class="success-icon">✓</div>
        <h4>Message Sent!</h4>
        <p>Thanks for reaching out. Our support team will get back to you at <strong>{{ form.email }}</strong> shortly.</p>
        <button class="btn btn-outline" @click="close">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';

const props = defineProps({
  isOpen: Boolean
});

const emit = defineEmits(['close']);

const isClosing = ref(false);
const submitted = ref(false);
const loading = ref(false);
const error = ref('');

const form = reactive({
  name: '',
  email: '',
  message: ''
});

const close = () => {
  isClosing.value = true;
  setTimeout(() => {
    emit('close');
    isClosing.value = false;
    if (submitted.value) {
        // reset after close
        setTimeout(() => {
            submitted.value = false;
            form.name = '';
            form.email = '';
            form.message = '';
        }, 300);
    }
  }, 300); // match animation duration
};

const submitForm = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
    const response = await fetch(`${API_URL}/support`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form)
    });
    
    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.error || 'Failed to send message');
    }
    
    submitted.value = true;
  } catch (err) {
    error.value = err.message || 'An error occurred. Please try again.';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(10, 10, 15, 0.7);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal-content {
  background: var(--bg-card, #1a1a24);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 16px;
  width: 90%;
  max-width: 450px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.modal-content.is-closing {
  animation: slideDown 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 1.5rem 1rem;
  border-bottom: 1px solid var(--border-color, rgba(255, 255, 255, 0.05));
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-primary, #fff);
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-secondary, #a0aec0);
  font-size: 1.5rem;
  cursor: pointer;
  transition: color 0.2s;
  line-height: 1;
}

.close-btn:hover {
  color: var(--text-primary, #fff);
}

.modal-body {
  padding: 1.5rem;
}

.subtitle {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: var(--text-secondary, #a0aec0);
  font-size: 0.95rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-secondary, #a0aec0);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem 1rem;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 8px;
  color: var(--text-primary, #fff);
  font-family: inherit;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color, #3b82f6);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.btn-block {
  width: 100%;
  margin-top: 1rem;
  padding: 0.875rem;
}

.error-message {
  color: #ef4444;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  text-align: center;
}

.success-state {
  text-align: center;
  padding: 2rem 1.5rem;
}

.success-icon {
  width: 60px;
  height: 60px;
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin: 0 auto 1.5rem;
}

.success-state h4 {
  margin: 0 0 0.5rem;
  color: var(--text-primary, #fff);
  font-size: 1.25rem;
}

.success-state p {
  color: var(--text-secondary, #a0aec0);
  margin-bottom: 2rem;
  font-size: 0.95rem;
  line-height: 1.5;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes slideDown {
  from { opacity: 1; transform: translateY(0) scale(1); }
  to { opacity: 0; transform: translateY(20px) scale(0.95); }
}
</style>
