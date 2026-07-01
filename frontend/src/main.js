import { createApp } from 'vue'
import './assets/css/common.css'
import 'font-awesome/css/font-awesome.min.css'
import App from './App.vue'
import router from './router'

import Toast from "vue-toastification"
import "vue-toastification/dist/index.css"

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const customDarkTheme = {
  dark: true,
  colors: {
    background: '#0b0f1a',
    surface: '#131929',
    primary: '#3b82f6',
    secondary: '#8b5cf6',
    error: '#f43f5e',
    info: '#2196F3',
    success: '#10b981',
    warning: '#eab308',
  }
}

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'customDarkTheme',
    themes: {
      customDarkTheme,
    }
  }
})

const app = createApp(App)

app.use(router)
app.use(Toast, {
    position: "top-right",
    timeout: 3000,
    closeOnClick: true,
    maxToasts: 3,
    newestOnTop: true,
    pauseOnFocusLoss: false,
})
app.use(vuetify)

app.mount('#app')
