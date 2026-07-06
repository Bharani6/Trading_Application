<template>
  <div class="modal-overlay" v-if="isOpen" @click.self="close">
    <div class="modal-content" v-bind:class="{ 'is-closing': isClosing }">
      <div class="modal-header">
        <h3>{{ title }}</h3>
        <button class="close-btn" @click="close">&times;</button>
      </div>

      <div class="modal-body">
        <div class="disclaimer-content">
          <p><strong>Risk Disclaimer:</strong> Investments in securities market are subject to market risks, read all the related documents carefully before investing. Brokerage will not exceed the SEBI prescribed limit.</p>
          <p>MoneymakePro acts as a technology platform and a registered broker. We do not provide investment advice or guarantee returns. Please ensure you understand the financial risks involved in trading before proceeding.</p>
          <template v-if="title === 'Terms of Service'">
            <TermsOfServiceContent />
          </template>
          <template v-if="title === 'Privacy Policy'">
            <PrivacyPolicyContent />
          </template>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-primary btn-block no-effect" @click="close">I Agree</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import TermsOfServiceContent from './TermsOfServiceContent.vue';
import PrivacyPolicyContent from './PrivacyPolicyContent.vue';

const props = defineProps({
  isOpen: Boolean,
  title: {
    type: String,
    default: 'Disclaimer'
  }
});

const emit = defineEmits(['close']);

const isClosing = ref(false);

const close = () => {
  isClosing.value = true;
  setTimeout(() => {
    emit('close');
    isClosing.value = false;
  }, 300); // match animation duration
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
  max-width: 800px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
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
  overflow-y: auto;
  max-height: 60vh;
}

.disclaimer-content {
  color: var(--text-secondary, #a0aec0);
  font-size: 0.95rem;
  line-height: 1.6;
}

:deep(.disclaimer-content p) {
  margin-bottom: 1rem;
}

:deep(.disclaimer-content strong) {
  color: var(--primary-color, #10b981);
}

:deep(.policy-section) {
  margin-bottom: 1.5rem;
}

:deep(.policy-section h4) {
  color: var(--text-primary, #fff);
  margin: 0 0 0.5rem 0;
  font-size: 1.05rem;
}

:deep(.policy-section ul) {
  padding-left: 1.5rem;
  margin: 0 0 1rem 0;
}

:deep(.policy-section li) {
  margin-bottom: 0.25rem;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid var(--border-color, rgba(255, 255, 255, 0.05));
}

.btn-block {
  width: 100%;
}

.no-effect {
  animation: none !important;
  transform: none !important;
}
.no-effect::before {
  animation: none !important;
  display: none !important;
}
.no-effect:hover {
  transform: none !important;
  box-shadow: none !important;
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
