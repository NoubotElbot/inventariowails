import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura' // ← cambio aquí
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import 'primeicons/primeicons.css'
import { router } from '@/router'
import App from './app/App.vue'
import './style.css'
import { definePreset } from '@primeuix/themes'

const Noir = definePreset(Aura, {
  semantic: {
    primary: {
      50: '{surface.50}',
      100: '{surface.100}',
      200: '{surface.200}',
      300: '{surface.300}',
      400: '{surface.400}',
      500: '{surface.500}',
      600: '{surface.600}',
      700: '{surface.700}',
      800: '{surface.800}',
      900: '{surface.900}',
      950: '{surface.950}',
      color: 'light-dark({primary.950}, {primary.50})',
      contrastColor: 'light-dark(#ffffff, {primary.950})',
      hoverColor: 'light-dark({primary.800}, {primary.200})',
      activeColor: 'light-dark({primary.700}, {primary.300})',
    },
    highlight: {
      background: 'light-dark({primary.950}, {primary.50})',
      focusBackground: 'light-dark({primary.700}, {primary.300})',
      color: 'light-dark(#ffffff, {primary.950})',
      focusColor: 'light-dark(#ffffff, {primary.950})',
    },
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(PrimeVue, {
  license: import.meta.env.VITE_PRIMEUI_LICENSE_KEY,
  theme: {
    preset: Noir,
    options: {
      darkModeSelector: '.dark',
    },
  },
})
app.use(ToastService)
app.use(ConfirmationService)

app.mount('#app')
